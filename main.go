package main

import (
	"log"

	"go-server/container"
	"go-server/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	c := container.NewContainer()

	routes.RegisterHealthRoutes(r)
	routes.RegisterUserRoutes(r, c)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}