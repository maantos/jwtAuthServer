package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/maantos/jwtAuth/handlers"
	"github.com/maantos/jwtAuth/initializers"
)

func init() {
	initializers.LoadEnvVarables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Group(func(r chi.Router) {
		r.Post("/signup", handlers.SignUp)
		r.Post("/login", handlers.Login)

	})

	// Private Routes
	// Require Authentication
	r.Group(func(r chi.Router) {
		r.Use(handlers.MyMiddleware)
		r.Get("/validate", handlers.Test)
	})

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
}
