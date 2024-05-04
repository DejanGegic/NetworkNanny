package db

import (
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type BadgerInstance struct {
	*badger.DB
}

func initBadger() *badger.DB {
	fmt.Println(os.Getenv("DB_LOCATION"))
	db, err := badger.Open(badger.DefaultOptions(os.Getenv("DB_LOCATION")).WithReadOnly(false))
	if err != nil {
		panic(err)
	}
	return db
}

func (b BadgerInstance) Write(key string, value string) error {
	txn := b.NewTransaction(true)
	if err := txn.Set([]byte(key), []byte(value)); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

func (b BadgerInstance) WriteTTL(key string, value string, ttl time.Duration) error {
	e := badger.NewEntry([]byte(key), []byte(value)).WithTTL(ttl)
	txn := b.NewTransaction(true)
	if err := txn.SetEntry(e); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}

	// save TTL
	keyTTL := key + "_ttl"
	//convert duration to time
	ttlTime := time.Now().Add(ttl)
	e = badger.NewEntry([]byte(keyTTL), []byte(ttlTime.Format(time.RFC3339))).WithTTL(ttl)
	txn = b.NewTransaction(true)
	if err := txn.SetEntry(e); err != nil {
		b.Delete(key)
		return err
	}
	if err := txn.Commit(); err != nil {
		b.Delete(key)
		return err
	}

	return nil
}

func (b BadgerInstance) Read(key string) (string, error) {
	var value string
	err := b.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
	})
	return value, err
}

func (b BadgerInstance) ReadTTL(key string) (value string, ttl time.Duration, err error) {

	// read data
	err = b.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
	})
	if err != nil {
		return value, ttl, err
	}

	// read ttl
	var ttlString string
	err = b.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key + "_ttl"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			ttlString = string(val)
			return nil
		})
	})
	if err != nil {
		return value, ttl, err
	}

	// convert string to time
	ttlTime, err := time.Parse(time.RFC3339, ttlString)
	if err != nil {
		return value, ttl, err
	}
	ttl = time.Since(ttlTime)
	return value, ttl, err

}

func (b BadgerInstance) Delete(key string) error {
	txn := b.NewTransaction(true)
	if err := txn.Delete([]byte(key)); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}
