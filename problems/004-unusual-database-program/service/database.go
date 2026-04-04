package service

import (
	"fmt"
	"log"
)

type Database map[string]string

func NewDatabase() Database {
	db := make(Database)
	db["version"] = "Sahil's Protocol"
	return db
}

func (db Database) Insert(key, value string) {
	if key == "version" {
		return
	}
	db[key] = value
	log.Println(key, value)
}

func (db Database) Retrieve(key string) string {
	value := db[key]
	return fmt.Sprintf("%s=%s", key, value)
}
