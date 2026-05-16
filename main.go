package main

import (
	"log"

	"go-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.RegisterHealthRoutes(r)
	routes.RegisterUserRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}