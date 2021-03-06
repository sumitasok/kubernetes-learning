package handler

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sumitasok/kubernetes-learning/internal"
	"github.com/sumitasok/kubernetes-learning/internal/store"
)

// FileStore holds the configs and connectors for interacting with the file system
// implements the handlers
type FileStore struct {
	MetaStore *store.InMemory
}

// AddFile adds the file to store
func (fS FileStore) AddFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]

	// TODO: When a single file out of multiple files upload fails, handle that scenario.
	for _, file := range files {
		log.Println("saving... ", file.Filename)

		tmpLocalPath := "/tmp/" + file.Filename
		localPath := "/store/" + file.Filename

		// TODO: verify that the uploaded file is txt file.

		// check if the file already exists in disk; redundant validation
		if _, err := os.Stat(localPath); err == nil {
			log.Println("file exists in disk ", localPath)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse("file with same name already exists - "+file.Filename))
			return
		}

		// savethe uploaded file to tmp location for md5 checksum
		if err := c.SaveUploadedFile(file, tmpLocalPath); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse(err.Error()))
			return
		}

		// find the md5 of the file.
		md5Value, err := internal.Md5(tmpLocalPath)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse(err.Error()))
			return
		}

		// save meta info to store; this will update us if a similar file with different name exists
		if err := fS.MetaStore.Add(file.Filename, md5Value); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse(err.Error()))
			return
		}

		// TODO: in case of a failure we need to revert the entry in metaStore.
		if err := internal.MoveFile(tmpLocalPath, localPath); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse(err.Error()))
			return
		}

		log.Println("saved ", file.Filename)
	}

	c.JSON(http.StatusAccepted, gin.H{"data": fS.MetaStore.Files, "status": "UPLOADED", "message": "file was uploaded succesfully"})
}

func failedFileUploadResponse(msg string) gin.H {
	return gin.H{"data": gin.H{}, "status": "FAILED", "message": errFailedUpload(msg)}
}

func errFailedUpload(msg string) string {
	return "file upload failed " + msg
}

// LsFile lists metainfo of the files in store
func (fS FileStore) LsFile(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{"data": fS.MetaStore.Files, "status": "DONE", "message": "successfully retrieved files"})
}

// UpdateFile Updates the file to store
func (fS FileStore) UpdateFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]

	// TODO: When a single file out of multiple files upload fails, handle that scenario.
	for _, file := range files {
		log.Println("saving... ", file.Filename)

		tmpLocalPath := "/tmp/" + file.Filename
		localPath := "/store/" + file.Filename

		// TODO: verify that the uploaded file is txt file.

		// save the uploaded file to tmp location for md5 checksum
		if err := c.SaveUploadedFile(file, tmpLocalPath); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse(err.Error()))
			return
		}

		// find the md5 of the file.
		md5Value, err := internal.Md5(tmpLocalPath)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse(err.Error()))
			return
		}

		// save meta info to store; this will update us if a similar file with different name exists
		// TODO: have separate method for Updating the files in MetaStore.
		if err := fS.MetaStore.Add(file.Filename, md5Value); err != nil {
			log.Println(err)
		}

		// TODO: in case of a failure we need to revert the entry in metaStore.
		if err := internal.MoveFile(tmpLocalPath, localPath); err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, failedFileUploadResponse(err.Error()))
			return
		}

		log.Println("saved ", file.Filename)
	}

	// TODO: distiguish between update and add.
	c.JSON(http.StatusAccepted, gin.H{"data": fS.MetaStore.Files, "status": "UPDATED", "message": "file was updated succesfully"})
}

// GetFile returns a single file to client
func (fS FileStore) GetFile(c *gin.Context) {
	targetPath := filepath.Join("/store/", c.Param("filename"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+c.Param("filename"))
	c.Header("Content-Type", "application/octet-stream")
	c.File(targetPath)
}

// DeleteFile delete a single file from server
func (fS FileStore) DeleteFile(c *gin.Context) {
	targetPath := filepath.Join("/store/", c.Param("filename"))

	f, err := fS.MetaStore.DeleteFileByName(c.Param("filename"))
	if err != nil {
		log.Println("Remove failed", err)
		c.JSON(http.StatusNotAcceptable, gin.H{"data": f, "status": "FAILED", "message": "file deletion failed " + err.Error()})
		return
	}

	if err := os.Remove(targetPath); err != nil {
		log.Println("Remove failed", err)
		c.JSON(http.StatusNotAcceptable, gin.H{"data": f, "status": "FAILED", "message": "file deletion failed " + err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"data": f, "status": "DONE", "message": "successfully deleted file"})
}
