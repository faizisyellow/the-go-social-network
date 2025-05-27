package main

import (
	"log"
	"os"

	"faizisyellow.github.com/thegosocialnetwork/internal/db"
	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := db.New(os.Getenv("DB_ADDRESS"), 30, 30, "15m")
	if err != nil {
		log.Panic(err)
	}

	defer conn.Close()

	log.Println("database connection pool established")

	db.Seed(store.NewStorage(conn))
}
