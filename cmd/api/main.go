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

	config := config{
		addr: os.Getenv("PORT"),
		db: dbConfig{
			addr:        os.Getenv("DB_ADDRESS"),
			maxOpenConn: 30,
			maxIdleConn: 30,
			maxIdleTime: "15m",
		},
	}

	db, err := db.New(config.db.addr, config.db.maxOpenConn, config.db.maxIdleConn, config.db.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	log.Println("database connection pool established")

	app := &application{
		config: config,
		store:  store.NewStorage(db),
	}

	log.Fatal(app.run(app.mount()))
}
