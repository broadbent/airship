package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

type DB struct {
	*bolt.DB
}

func (db *DB) Open(path string, mode os.FileMode, buckets []string) error {
	var err error
	db.DB, err = bolt.Open(path, mode, nil)
	if err != nil {
		return err
	}

	for _, bucket := range buckets {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
			return nil
		})
	}

	return nil
}

func (db *DB) Write(bucket string, key string, value string, timestamp bool) string {
	if timestamp {
		key = time.Now().UTC().Format(time.RFC3339)
	}

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		err := bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return key
}

func (db *DB) Read(bucket string, key string) {

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucket))

		val := bucket.Get([]byte(key))
		fmt.Println(string(val))

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
