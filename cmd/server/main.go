package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/sumitasok/kubernetes-learning/internal/handler"
	"github.com/sumitasok/kubernetes-learning/internal/store"
)

func main() {
	log.Println("server...")

	router := gin.Default()

	// TODO: factory pattern
	fileStore := handler.FileStore{
		MetaStore: store.NewInMemory(),
	}

	router.POST("/files", fileStore.AddFile)

	router.Run(os.Getenv("PORT"))
	// router.Run(":8080")
}
