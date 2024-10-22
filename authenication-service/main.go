package main

import (
	"authentication-service/config"
	"authentication-service/handler"
	"authentication-service/middlerware"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	db := config.ConnectDB()
	config.UpdateDB(db)
	r := chi.NewRouter()
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

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Println("Server started at port", port)
	http.ListenAndServe(port, r)
}
