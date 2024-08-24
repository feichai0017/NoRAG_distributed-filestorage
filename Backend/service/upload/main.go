package main

import (
	"cloud_distributed_storage/handler"
	"fmt"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	// Create a new CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://192.168.0.200:3001"},              // 明确指定允许的源
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"}, // 添加OPTIONS
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})

	// Create a new ServeMux
	mux := http.NewServeMux()

	mux.HandleFunc("/api/upload", handler.HTTPInterceptor(handler.UploadHandler))
	mux.HandleFunc("/api/upload/success", handler.HTTPInterceptor(handler.UploadSucHandler))
	mux.HandleFunc("/api/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))
	mux.HandleFunc("/api/meta/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	mux.HandleFunc("/api/download", handler.HTTPInterceptor(handler.DownloadHandler))
	mux.HandleFunc("/api/update", handler.HTTPInterceptor(handler.FileMetaUpdateHandler))
	mux.HandleFunc("/api/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))
	mux.HandleFunc("/api/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	//redis upload api
	mux.HandleFunc("/api/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	mux.HandleFunc("/api/mpupload/uploadpart", handler.HTTPInterceptor(handler.UploadPartHandler))
	mux.HandleFunc("/api/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))

	//user api
	mux.HandleFunc("/api/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	mux.HandleFunc("/api/login", handler.SignInHandler)
	mux.HandleFunc("/api/signup", handler.SignupHandler)

	// Wrap the mux with the CORS handler
	handler := c.Handler(mux)

	// Start the server
	err := http.ListenAndServe(":8081", handler)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s", err.Error())
	}
}
