package model

import (
	"github.com/boltdb/bolt"
	"io/ioutil"
	"os"
	"path"
	"log"
)

var db *bolt.DB

func init() {
	dbPath, err := ioutil.TempFile("/tmp", "glinks")
	if err != nil {
		log.Panic("Could not create tempfile")
	}
	log.Printf("DB path: %s", dbPath)
	os.MkdirAll(path.Dir(dbPath), 0755)
	db, _ = bolt.Open(dbPath, 0644, nil)
}

func GetDB() *bolt.DB {
	return db
}
