package main

import (
	"bytes"
	"encoding/json"
	"io"
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
	switch os.Args[1] {
	case "add":
		add(os.Args[2:]...)
	case "ls":
		ls()
	}
}

// FileResp holds the response from files
type FileResp struct {
	Data map[string]struct {
		Checksum string `json:'Checksum'`
		name     string `json:'Name'`
	} `json:'data'`
	Status  string `json:'status'`
	Message string `json:'message'`
}

func ls() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", RemoteStoreBaseURL+"/files", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("coun't list files ", err.Error())
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var t FileResp
	err = decoder.Decode(&t)

	if err != nil {
		log.Println("coun't list files ", err.Error())
	}

	for k := range t.Data {
		print(k, "\n")
	}
}

func add(files ...string) {
	client := &http.Client{}
	for _, _filepath := range files {
		file, err := os.Open(_filepath)
		if err != nil {
			logCouldntAddFile(_filepath, err.Error())
			continue
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", filepath.Base(_filepath))
		if err != nil {
			logCouldntAddFile(_filepath, err.Error())
			continue
		}
		_, err = io.Copy(part, file)

		err = writer.Close()
		if err != nil {
			logCouldntAddFile(_filepath, err.Error())
			continue
		}

		req, err := http.NewRequest("POST", RemoteStoreBaseURL+"/files", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := client.Do(req)
		if err != nil {
			logCouldntAddFile(_filepath, err.Error())
			continue
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		var t FileResp
		err = decoder.Decode(&t)

		log.Println("add file: ", _filepath, t.Status, t.Message)
	}
}

func logCouldntAddFile(filename, message string) {
	log.Println("coun't add file ", filename, message)
}
