package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/maantos/jwtAuth/pkg/db"
	"github.com/maantos/jwtAuth/pkg/handlers"
)

func init() {
	//initializers.LoadEnvVarables()
	db.ConnectToDB()
	db.SyncDatabase()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Post("/signup", handlers.SignUp)
		r.Post("/login", handlers.Login)
		r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("service is up."))
		})

	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(handlers.MyMiddleware)
		r.Get("/validate", handlers.Test)
	})
	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	l.Printf("Starting server on port %s...", os.Getenv("PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
}
