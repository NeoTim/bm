package main

import (
	"bytes"
	"github.com/boltdb/bolt"
	"log"
	"path"
)

var db *bolt.DB

const (
	DefaultDBName = "bm.db"
	DefaultBucket = "bm"
)

// OpenDB opens the database.
func OpenDB(config Config) {
	var err error
	db, err = bolt.Open(path.Join(config.StorePath, DefaultDBName), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

}

// CloseDB closes the database.
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// InitDBBucket creates the bucket if it doesn't exit.
func InitDBBucket() {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DefaultBucket))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
}

// Get the value associated with the input key.
func Get(key string) string {
	var value []byte = nil
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DefaultBucket))
		value = b.Get([]byte(key))
		return nil
	})
	return string(value)
}

// Put the (key, value) pair into database.
func Put(key string, value string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DefaultBucket))
		err := b.Put([]byte(key), []byte(value))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
}

// Delete the (key, value) pair from database.
func Delete(key string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DefaultBucket))
		err := b.Delete([]byte(key))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
}

// IterateKey returns the key set with a prefix.
func IterateKey(prefix string) []string {
	var results []string

	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(DefaultBucket)).Cursor()

		pre := []byte(prefix)
		for k, _ := c.Seek(pre); k != nil && bytes.HasPrefix(k, pre); k, _ = c.Next() {
			results = append(results, string(k))
		}
		return nil
	})
	return results
}
