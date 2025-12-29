package main

import (
	"fmt"
	"log"

	"go1f/pkg/api"
	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	err := db.Init("")
	if err != nil {
		log.Fatalf("Не удалось инициализировать базу данных: %v", err)
	}
	fmt.Println("Database initialized successfully")

	api.Init()

	server.Start()
}
