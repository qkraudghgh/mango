package mangoUtils

import (
	"encoding/binary"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/qkraudghgh/mango/manager"
)

// GetDbPath gets dbPath corresponding to each OS.
func GetDbPath() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

		return home
	}

	dbPath := filepath.Join(os.Getenv("HOME"), ".mango.db")

	return dbPath
}

// CheckBucketAndMake function check bucket, and If not exist make bucket
func CheckBucketAndMake() {
	db, err := bolt.Open(GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(manager.MangoBucket))
		if err != nil {
			return err
		}

		return nil
	})
}

// CheckKey function is stored and return an error if not
func CheckKey(key int) error {
	db, err := bolt.Open(GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(manager.MangoBucket))
		v := b.Get(Itob(key))
		if v == nil {
			return errors.New("That todo does not exist")
		}
		return nil
	})

	return err
}

// Itob function convert integer to byte
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// ValidateArgs function check Arguments length and type
func ValidateArgs(args []string) (int, error) {
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
