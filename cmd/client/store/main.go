package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// RemoteStoreBaseURL the address of the remote server
// TODO: move this to aconfig accessible from env
const RemoteStoreBaseURL = "http://localhost:8080"

func main() {
	print("hi client\n")

	switch os.Args[1] {
	case "add":
		add(os.Args[2:]...)
	}
}

func add(files ...string) {
	client := &http.Client{}
	for _, _filepath := range files {
		file, err := os.Open(_filepath)
		if err != nil {
			log.Println("coun't add file ", file, err.Error())
			continue
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(_filepath))
		if err != nil {
			log.Println("coun't add file ", file, err.Error())
			continue
		}
		_, err = io.Copy(part, file)

		err = writer.Close()
		if err != nil {
			log.Println("coun't add file ", file, err.Error())
			continue
		}

		req, err := http.NewRequest("POST", RemoteStoreBaseURL+"/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := client.Do(req)
		if err != nil {
			log.Println("coun't add file ", file, err.Error())
			continue
		}
		defer resp.Body.Close()

		_body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("coun't add file ", file, err.Error())
			continue
		}
		log.Println("added file ", file, string(_body))
	}
}
