package main

import (
	"go-server/container"
	"go-server/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	c, err := container.NewContainer()
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("failed to close container: %v", err)
		}
	}()

	routes.RegisterHealthRoutes(r)
	routes.RegisterUserRoutes(r, c)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
