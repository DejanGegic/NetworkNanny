package db

import (
	"time"

	"github.com/dgraph-io/badger/v4"
)

var db = initDb()

func initDb() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("./badger").WithReadOnly(false))
	if err != nil {
		panic(err)
	}
	return db
}
func Write(key string, value string) error {
	txn := db.NewTransaction(true)
	if err := txn.Set([]byte(key), []byte(value)); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

func WriteTTL(key string, value string, ttl time.Duration) error {
	e := badger.NewEntry([]byte(key), []byte(value)).WithTTL(ttl)
	txn := db.NewTransaction(true)
	if err := txn.SetEntry(e); err != nil {
		return err
	}
	if err := txn.Commit(); err != nil {
		return err
	}
	return nil
}

func Read(key string) (string, error) {
	var value string
	err := db.View(func(txn *badger.Txn) error {
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
