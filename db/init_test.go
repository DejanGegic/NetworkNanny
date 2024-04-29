package db

import (
	"os"
	"reflect"
	"testing"
)

func TestInitDB(t *testing.T) {

	os.Setenv("DB_TYPE", "redis")
	db := InitDB()
	if reflect.TypeOf(db).String() != "*db.RedisInstance" {
		t.Errorf("InitDB() = %v, want %v", reflect.TypeOf(db).String(), "*db.RedisInstance")
	}

	db.Write("test", "test")

	result, err := db.Read("test")
	if err != nil {
		t.Errorf("Read() error = %v", err)
	}
	if result != "test" {
		t.Errorf("Read() = %v, want %v", result, "test")
	}

	// os.Setenv("DB_TYPE", "badger")
	// db = InitDB()
	// if reflect.TypeOf(db).String() != "*db.BadgerInstance" {
	// 	t.Errorf("InitDB() = %v, want %v", reflect.TypeOf(db).String(), "*db.BadgerInstance")
	// }
}
