package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"encoding/binary"
	"github.com/boltdb/bolt"
	"github.com/qkraudghgh/mango/manager"
	"log"
	"strconv"
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
	_, err := validateArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	db, err := bolt.Open(manager.GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
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
	todoNo, err := validateArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	manager.CheckBucketAndMake()

	db, err := bolt.Open(manager.GetDbPath(), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = checkKey(todoNo)
	if err != nil {
		return err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(manager.MangoBucket))
		err := b.Delete(itob(todoNo))
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

	manager.CheckBucketAndMake()

	db, err := bolt.Open(manager.GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
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
	todoNo, err := validateArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	manager.CheckBucketAndMake()

	err = checkKey(todoNo)
	if err != nil {
		return err
	}

	updateIsChecked(todoNo, 1)

	return nil
}

func unDoneFunc(args []string) error {
	todoNo, err := validateArgs(args)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	manager.CheckBucketAndMake()

	err = checkKey(todoNo)
	if err != nil {
		return err
	}

	updateIsChecked(todoNo, 0)

	return nil
}

func updateIsChecked(keyId int, isCheck int) error {
	db, err := bolt.Open(manager.GetDbPath(), 0755, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var oldTodo Todo

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(manager.MangoBucket))
		v := b.Get(itob(keyId))

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
		return b.Put(itob(keyId), encoded)
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
		return b.Put(itob(todo.ID), encoded)
	})
	return err
}

func checkKey(key int) error {
	db, err := bolt.Open(manager.GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(manager.MangoBucket))
		v := b.Get(itob(key))
		if v == nil {
			return errors.New("That todo does not exist")
		}
		return nil
	})

	return err
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func validateArgs(args []string) (int, error) {
	nArgs := len(args)

	var todoNo int
	var err error

	if nArgs != 1 {
		err = errors.New("Invalid arguments: this command could take one argument at most")
	} else {
		if todoNo, err = strconv.Atoi(args[0]); err != nil {
			err = errors.New("Integer is allowed only")
		}
	}

	return todoNo, err
}