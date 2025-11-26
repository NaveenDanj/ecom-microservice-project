package main

import (
	"fmt"
	"net/http"
	"user-service/config"
	"user-service/models"
	"user-service/routes"
	"user-service/utils"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{})

	_, err := utils.InitRabbitMQ("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic("Could not connect to RabbitMQ: " + err.Error())
	}

	r := mux.NewRouter()
	routes.UserRoutes(r)

	fmt.Println("User service running on port 8081")
	http.ListenAndServe(":8081", r)
}
