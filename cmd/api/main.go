package main

import (
	"log"
	"os"
	"time"

	"faizisyellow.github.com/thegosocialnetwork/internal/db"
	"faizisyellow.github.com/thegosocialnetwork/internal/helpers"
	"faizisyellow.github.com/thegosocialnetwork/internal/mailer"
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

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: better config with default value (to do debug need default key)
	config := config{
		addr: helpers.DefaultString(os.Getenv("PORT"), ":8080"),
		db: dbConfig{
			addr:        os.Getenv("DB_ADDRESS"),
			maxOpenConn: 30,
			maxIdleConn: 30,
			maxIdleTime: "15m",
		},
		env: helpers.DefaultString(os.Getenv("ENV"), "Development"),
		mail: mailConfig{
			fromEmail: helpers.DefaultString(os.Getenv("FROM_EMAIL"), ""),
			sendGrid: sendgridConfig{
				apiKey: helpers.DefaultString(os.Getenv("SENDGRID_API_KEY"), ""),
			},
			exp: time.Hour * 24 * 3, // 3 days
		},
		frontendURL: helpers.DefaultString(os.Getenv("FRONTEND_URL"), "http://localhost:4173"),
		auth: authConfig{
			basic: basicAuthConfig{
				user: helpers.DefaultString(os.Getenv("AUTH_BASIC_USER"), "admin"),
				pass: helpers.DefaultString(os.Getenv("AUTH_BASIC_PASSWORD"), "admin"),
			},
		},
	}

	//TODO: fix the error logger in error.go
	logger := zap.Must(zap.NewProduction()).Sugar()

	defer logger.Sync()

	db, err := db.New(config.db.addr, config.db.maxOpenConn, config.db.maxIdleConn, config.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	mailer := mailer.NewSendgrid(config.mail.sendGrid.apiKey, config.mail.fromEmail)

	app := &application{
		config: config,
		store:  store.NewStorage(db),
		logger: logger,
		mailer: mailer,
	}

	logger.Fatal(app.run(app.mount()))
}
