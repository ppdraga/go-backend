package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

func main() {
	uploadHandler := &UploadHandler{
		HostAddr:  "localhost:8081",
		UploadDir: "upload",
	}
	http.Handle("/upload", uploadHandler)

	// FileServer in background
	go func() {
		dirToServe := http.Dir(uploadHandler.UploadDir)

		fs := &http.Server{
			Addr:         ":8081",
			Handler:      http.FileServer(dirToServe),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		err := fs.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	// Upload Server
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Println(err)
	}

}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}
	filePath := h.UploadDir + "/" + header.Filename

	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fileLink := h.HostAddr + "/" + header.Filename
	//fmt.Fprintf(w, "File %s has been successfully uploaded", header.Filename)

	// Checking file
	req, err := http.NewRequest(http.MethodHead, fileLink, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to check file", http.StatusInternalServerError)
		return
	}
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to check file", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, fileLink)
}
