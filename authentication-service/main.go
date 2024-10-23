package main

import (
	"authentication-service/config"
	"authentication-service/handler"
	"authentication-service/middlerware"
	"fmt"
	"log"
	"net/http"

	_ "authentication-service/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	db := config.ConnectDB()
	config.UpdateDB(db)
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	user := handler.NewUserHandler()
	r.Group(func(router chi.Router) {
		router.Use(middlerware.Authorize)
		router.Get("/users", user.GetUsers)
		router.Get("/users/{email}", user.GetUserByEmail)
		router.Get("/users/{id}", user.GetUserById)

	})
	// Public routes (no auth middleware)
	r.Post("/users/register", user.CreateUser)
	r.Post("/login", user.Login)
	r.Get("/authorize", user.Authentication)

	fmt.Println("Server started at port :3008")
	if err := http.ListenAndServe(":3008", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
