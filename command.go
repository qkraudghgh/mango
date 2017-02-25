package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/qkraudghgh/mango/manager"
	"github.com/qkraudghgh/mango/utils"
	"os"
)

var (
	addCommand = manager.Command{
		Name:  "add",
		Usage: "add    : mango add 'todo'",
		Run:   addFunc,
	}
	listCommand = manager.Command{
		Name:  "list",
		Usage: "list   : mango list",
		Run:   listFunc,
	}
	deleteCommand = manager.Command{
		Name:  "delete",
		Usage: "delete : mango delete [number]",
		Run:   deleteFunc,
	}
	doneCommand = manager.Command{
		Name:  "done",
		Usage: "done : mango done [number]",
		Run:   doneFunc,
	}
	unDoneCommand = manager.Command{
		Name:  "undone",
		Usage: "undone : mango undone [number]",
		Run:   unDoneFunc,
	}
	clearCommand = manager.Command{
		Name:  "clear",
		Usage: "clear : mango clear",
		Run:   clearFunc,
	}
)

// mango add 'todo'
func addFunc(args []string) error {
	if len(args) != 1 {
		return errors.New("add command needs only one argument")
	}

	// connect to DB
	db, err := bolt.Open(mangoUtils.GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// make structure TODO
	newTodo := Todo{
		Content:   args[0],
		CreatedAt: time.Now(),
		IsCheck:   0,
	}

	// save Todo
	if err := newTodo.save(db); err != nil {
		return errors.New("Error Write the todo to file")
	}

	return nil
}

// mango delete [number]
func deleteFunc(args []string) error {
	todoNo, err := mangoUtils.ValidateArgs(args)
	if err != nil {
		return err
	}

	// If not exist bucket, make bucket
	if err := mangoUtils.CheckBucketAndMake(); err != nil {
		return err
	}

	// check data from key
	if err := mangoUtils.CheckKey(todoNo); err != nil {
		return err
	}

	// connect to DB
	db, err := bolt.Open(mangoUtils.GetDbPath(), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// delete data
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(manager.MangoBucket))
		err := b.Delete(mangoUtils.Itob(todoNo))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errors.New("Error Delete the todo to file")
	}

	fmt.Printf("Todo #%d was deleted\n", todoNo)

	return nil
}

// mango list
func listFunc(args []string) error {
	if len(args) != 0 {
		return errors.New("The list command do not use argument")
	}

	// If not exist bucket, make bucket
	if err := mangoUtils.CheckBucketAndMake(); err != nil {
		return err
	}

	// connect to DB
	db, err := bolt.Open(mangoUtils.GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	todos := []Todo{}

	// make Json object array [{Id: 1, IsCheck: 1, Content: "blah blah"}, ...]
	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(manager.MangoBucket))

		b.ForEach(func(k, v []byte) error {
			var todo Todo

			if err := json.Unmarshal(v, &todo); err != nil {
				return err
			}

			todos = append(todos, todo)
			return nil
		})
		return nil
	})

	// print todos
	if err := PrintTodos(todos); err != nil {
		return err
	}

	return err
}

// mango done [number]
func doneFunc(args []string) error {

	// validate arguments
	todoNo, err := mangoUtils.ValidateArgs(args)
	if err != nil {
		return err
	}

	// If not exist bucket, make bucket
	if err := mangoUtils.CheckBucketAndMake(); err != nil {
		return err
	}

	// check data from key
	if err := mangoUtils.CheckKey(todoNo); err != nil {
		return err
	}

	// update isCheck to true
	if err := updateIsChecked(todoNo, 1); err != nil {
		return err
	}

	fmt.Printf("Todo #%d was done\n", todoNo)

	return nil
}

// mango undone [number]
func unDoneFunc(args []string) error {

	// validate arguments
	todoNo, err := mangoUtils.ValidateArgs(args)
	if err != nil {
		return err
	}

	// If not exist bucket, make bucket
	if err := mangoUtils.CheckBucketAndMake(); err != nil {
		return err
	}

	// check data from key
	if err := mangoUtils.CheckKey(todoNo); err != nil {
		return err
	}

	// update isCheck to false
	if err := updateIsChecked(todoNo, 0); err != nil {
		return err
	}

	fmt.Printf("Todo #%d was undone\n", todoNo)

	return nil
}

// mango clear
func clearFunc(args []string) error {
	// check arguments
	if len(args) != 0 {
		return errors.New("The clear command do not use argument")
	}

	// delete Db file
	os.Remove(mangoUtils.GetDbPath())

	fmt.Println("Todo is clear")

	return nil
}

// this function update 'isCheck' value from key
func updateIsChecked(keyID int, isCheck int) error {

	// connect to DB
	db, err := bolt.Open(mangoUtils.GetDbPath(), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var oldTodo Todo

	// store data to oldTodo
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(manager.MangoBucket))
		v := b.Get(mangoUtils.Itob(keyID))

		json.Unmarshal(v, &oldTodo)
		return nil
	})

	// modify isCheck value
	oldTodo.IsCheck = isCheck

	// update todo value
	db.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(oldTodo)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(manager.MangoBucket))
		return b.Put(mangoUtils.Itob(keyID), encoded)
	})

	return nil
}

// save todo method
func (todo *Todo) save(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		// If not exist bucket, make bucket
		b, err := tx.CreateBucketIfNotExists([]byte(manager.MangoBucket))
		if err != nil {
			return err
		}

		// auto increment id (Primary key)
		id, _ := b.NextSequence()
		todo.ID = int(id)

		encoded, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		return b.Put(mangoUtils.Itob(todo.ID), encoded)
	})
	return err
}
