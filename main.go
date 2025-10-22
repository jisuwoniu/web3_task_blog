package main

import (
	"log"
	"web3_task_blog/internal/routes"
)

func main() {

	// /HTTP
	r := routes.SetupRoutes()
	log.Println("Server is running on :8080")
	log.Fatal(r.Run(":8080"))
}
