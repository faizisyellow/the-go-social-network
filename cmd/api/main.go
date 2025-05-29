package main

import (
	"os"

	"faizisyellow.github.com/thegosocialnetwork/internal/db"
	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			The Go Social Network Restful API
//	@version		1.0
//	@description	Restful Api for The Go Social Network, the media social for ghopers.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

//

// @securityDefinitions.apikey	ApikeyAuth
// @in							header
// @name						Authorization
// @decsription
func main() {

	//TODO: fix the error logger in error.go
	loggerzap, _ := zap.NewProduction()
	defer loggerzap.Sync()

	logger := loggerzap.Sugar()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	// TODO: better config with default value (to do debug need default key)
	config := config{
		addr: os.Getenv("PORT"),
		db: dbConfig{
			addr:        os.Getenv("DB_ADDRESS"),
			maxOpenConn: 30,
			maxIdleConn: 30,
			maxIdleTime: "15m",
		},
		env: os.Getenv("ENV"),
	}

	db, err := db.New(config.db.addr, config.db.maxOpenConn, config.db.maxIdleConn, config.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: config,
		store:  store.NewStorage(db),
		logger: logger,
	}

	logger.Fatal(app.run(app.mount()))
}
