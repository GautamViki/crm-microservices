package main

import (
	"crm-service/config"
	"crm-service/controller"
	"crm-service/middlerware"
	"crm-service/models"
	"fmt"
	"log"
	"net/http"

	_ "crm-service/docs" // The generated docs package

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	config := config.ConnectDB()
	models.Migrate(config.DB)
	controllers := controller.NewHandler()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Protect routes with auth middleware
	auth := r.Group("/")              // You can change this group to something more specific if needed
	auth.Use(middlerware.Authorize()) // Apply the middleware

	// Routes for CRUD operations
	auth.POST("/customers/upload", controllers.UploadCustomer)
	auth.GET("/customers", controllers.ListCustomers)
	auth.PUT("/customers/:id", controllers.UpdateCustomer)
	auth.DELETE("/customers/:id", controllers.DeleteCustomer)
	auth.GET("/customers/cache", controllers.GetAllCacheCustomers)
	auth.GET("/customers/:id", controllers.GetUserById)

	fmt.Println("Server running at http://localhost:3007")
	r.Run(":3007")
}
