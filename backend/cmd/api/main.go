package main

import (
	"complaints/cmd/api/models"
	"complaints/cmd/api/routes"
	"log"
	"net/http"
)

func main() {

	err := models.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	router := routes.InitRoutes() // Call the InitRoutes function
	port := "4000"
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
