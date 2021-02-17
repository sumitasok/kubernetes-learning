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

	// TODO: every time server comes up, read the files and create checksum and save it in memory.
	// While inmemory store simplifies the checksum creation by caching the same

	router.POST("/files", fileStore.AddFile)
	// Ideally a file identifier should be at the end of the URL
	router.PUT("/files", fileStore.UpdateFile)
	router.GET("/files", fileStore.LsFile)
	router.GET("/files/:filename", fileStore.GetFile)

	router.Run(os.Getenv("PORT"))
}
