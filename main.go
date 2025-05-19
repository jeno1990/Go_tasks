package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("couldn’t load config: %v", err)
	}

	db, err := NewMysqlStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initStorage(db)

	err = createTable(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating table: %v\n", err)
		os.Exit(1)
	}
	handleCommand(db)
}
