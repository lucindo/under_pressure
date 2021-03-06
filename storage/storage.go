// Package storage provides helper functions to read/write pressure data points
package storage

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/lucindo/under_pressure/log"
	"github.com/lucindo/under_pressure/pressure"
)

// private boltdb vars
var (
	db     *bolt.DB
	bucket = "pressure"
)

// Init initializes the database
func Init(dbFile string) {
	var err error

	log.Logger.Printf("openning database file: %s\n", dbFile)
	db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Logger.Printf("error opening database file [%s]: %v\n", dbFile, err)
		panic(err)
	}

	log.Logger.Printf("creating bucket [%s] (if not exists)\n", dbFile)
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Logger.Printf("error creating bucket [%s]: %v\n", bucket, err)
		panic(err)
	}
}

// Close the databse
func Close() {
	log.Logger.Printf("closing database\n")
	err := db.Close()
	if err != nil {
		log.Logger.Printf("error closing database: %v", err)
	}
}

// AddPressure adds a new pressure datapoint to the database
func AddPressure(p pressure.Pressure) error {
	log.Logger.Printf("adding new pressure point to the database: %s\n", p)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		encoded, err := json.Marshal(p)
		if err != nil {
			return err
		}
		return b.Put([]byte(fmt.Sprintf("%d", p.Timestamp)), encoded)
	})
	if err != nil {
		log.Logger.Printf("error adding %s to database: %v\n", p, err)
		return err
	}
	return nil
}

// ListPressures lists all pressure points on database
func ListPressures() ([]pressure.Pressure, error) {
	var pList []pressure.Pressure
	log.Logger.Printf("getting all pressure points")
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.ForEach(func(k, v []byte) error {
			var decoded pressure.Pressure
			err := json.Unmarshal(v, &decoded)
			if err != nil {
				log.Logger.Printf("error decoding pressure point key [%v]: %v\n", k, err)
				return err
			}
			pList = append(pList, decoded)
			return nil
		})
		return err
	})
	if err != nil {
		log.Logger.Printf("error reading database: %v\n", err)
		return nil, err
	}
	log.Logger.Printf("got %d pressure points from database\n", len(pList))
	return pList, nil
}

// DeletePressure removes a pressure point by key (timestamp)
func DeletePressure(key string) error {
	log.Logger.Printf("removing key: %s", key)
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete([]byte(key))
		return err
	})
	if err != nil {
		log.Logger.Printf("error removing key[%s]: %v\n", key, err)
	}
	return err
}
