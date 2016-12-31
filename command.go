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
)

func addFunc(args []string) error {
	if len(args) != 1 {
		return errors.New("add command needs only one argument")
	}

	db, err := bolt.Open(mangoUtils.GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	newTodo := Todo{
		Content:   args[0],
		CreatedAt: time.Now(),
		IsCheck:   0,
	}

	err = newTodo.save(db)
	if err != nil {
		return errors.New("Error Write the todo to file")
	}

	return nil
}

func deleteFunc(args []string) error {
	todoNo, err := mangoUtils.ValidateArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	mangoUtils.CheckBucketAndMake()

	db, err := bolt.Open(mangoUtils.GetDbPath(), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = mangoUtils.CheckKey(todoNo)
	if err != nil {
		return err
	}

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

func listFunc(args []string) error {
	if len(args) != 0 {
		return errors.New("The list command do not use argument")
	}

	mangoUtils.CheckBucketAndMake()

	db, err := bolt.Open(mangoUtils.GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	todos := []Todo{}

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(manager.MangoBucket))

		b.ForEach(func(k, v []byte) error {
			var todo Todo
			err := json.Unmarshal(v, &todo)
			if err != nil {
				return err
			}
			todos = append(todos, todo)
			return nil
		})
		return nil
	})

	PrintTodos(todos)

	return nil
}

func doneFunc(args []string) error {
	todoNo, err := mangoUtils.ValidateArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	mangoUtils.CheckBucketAndMake()

	err = mangoUtils.CheckKey(todoNo)
	if err != nil {
		return err
	}

	updateIsChecked(todoNo, 1)

	fmt.Printf("Todo #%d was done\n", todoNo)

	return nil
}

func unDoneFunc(args []string) error {
	todoNo, err := mangoUtils.ValidateArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	mangoUtils.CheckBucketAndMake()

	err = mangoUtils.CheckKey(todoNo)
	if err != nil {
		return err
	}

	updateIsChecked(todoNo, 0)

	fmt.Printf("Todo #%d was undone\n", todoNo)

	return nil
}

func updateIsChecked(keyId int, isCheck int) error {
	db, err := bolt.Open(mangoUtils.GetDbPath(), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var oldTodo Todo

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(manager.MangoBucket))
		v := b.Get(mangoUtils.Itob(keyId))

		json.Unmarshal(v, &oldTodo)
		return nil
	})

	oldTodo.IsCheck = isCheck

	db.Update(func(tx *bolt.Tx) error {
		encoded, err := json.Marshal(oldTodo)
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(manager.MangoBucket))
		return b.Put(mangoUtils.Itob(keyId), encoded)
	})

	return nil
}

func (todo *Todo) save(db *bolt.DB) error {
	// Store the user model in the user bucket using the username as the key.
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(manager.MangoBucket))
		if err != nil {
			return err
		}

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
