package main

import (
	"crm-backend/config"
	"crm-backend/controller"
	"crm-backend/models"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	_ "crm-backend/docs"

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
	// Routes for CRUD

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/customers/upload", controllers.UploadCustomer)
	r.GET("/customers", controllers.ListCustomers)
	r.PUT("/customers/:id", controllers.UpdateCustomer)
	r.DELETE("/customers/:id", controllers.DeleteCustomer)
	r.GET("/customers/cache", controllers.GetAllCacheCustomers)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	fmt.Println("Server running at http://localhost" + port)
	r.Run(port)
}
