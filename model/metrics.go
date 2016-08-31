package model

import (
	"github.com/mmessmore/glinks/glinks"
	"github.com/boltdb/bolt"
	"fmt"
	"log"
	"bytes"
	"encoding/gob"
)

const bucket string = "nunchucks"

func Get(key string, db bolt.DB) (glinks.Metric){
	var value glinks.Metric
	var val bytes.Buffer
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}
		&val = bucket.Get(key)
		return nil
	})
	if err != nil {
		log.Panic("Could not find key")
	}
	dec := gob.NewDecoder(&value)
	dec.Decode(&val)
	return value
}

func Set(key string, val glinks.Metric, db bolt.DB) (error) {
	value := &bytes.Buffer{}
	enc := gob.NewEncoder(&value)
	err := enc.Encode(&val)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return err
		}

		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}
