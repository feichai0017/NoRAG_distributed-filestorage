package main

import (
	"cloud_distributed_storage/Backend/route"
	"fmt"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	// Initialize Gin router
	r := route.Router()

	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://192.168.0.200:3001"},              // 明确指定允许的源
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"}, // 添加OPTIONS
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	// Create a new ServeMux
	// Wrap the mux with the CORS handler
	handler := c.Handler(r)

	// Start the server
	err := http.ListenAndServe(":8081", handler)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}
