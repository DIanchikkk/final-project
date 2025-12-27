package main

import (
	"fmt"
	"log"
	"os"

	"go1f/pkg/api"
	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	// Получаем путь к файлу базы данных из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем базу данных
	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Не удалось инициализировать базу данных: %v", err)
	}
	fmt.Println("Database initialized successfully")
	api.Init()

	// Запускаем веб-сервер
	server.Start()
}
