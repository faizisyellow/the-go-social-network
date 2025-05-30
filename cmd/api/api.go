package main

import (
	"fmt"
	"net/http"
	"time"

	"faizisyellow.github.com/thegosocialnetwork/docs"
	"faizisyellow.github.com/thegosocialnetwork/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type dbConfig struct {
	addr        string
	maxOpenConn int
	maxIdleConn int
	maxIdleTime string
}

type mail struct {
	exp time.Duration
}

type config struct {
	addr string
	db   dbConfig
	env  string
	mail
}

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// proccesing should be stop
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.PostContextMiddleware)

				r.Get("/", app.getPostHandler)
				r.Patch("/", app.UpdatePostHandler)
				r.Delete("/", app.DeletePostHandler)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", app.createCommentHandler)
				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)

				// userID is the ID of the user we want to follow.
				r.Put("/follow/", app.followUserHandler)
				r.Put("/unfollow/", app.unFollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})

			r.Route("/authentication", func(r chi.Router) {
				r.Get("/", app.registerUserHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost%v", app.config.addr)

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server has started", "addr", fmt.Sprintf("localhost%v", app.config.addr), "env", app.config.env)

	return srv.ListenAndServe()
}
