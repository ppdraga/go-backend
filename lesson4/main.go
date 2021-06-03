package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type UploadHandler struct {
	HostAddr  string
	UploadDir string
}

type BrowseHandles struct {
	BrowseDir string
}

type FileInfo struct {
	Filename  string `json:"filename"`
	Extension string `json:"extension"`
	Size      int64  `json:"size"`
}

var FILE_EXT_REGEXP *regexp.Regexp

func init() {
	FILE_EXT_REGEXP = regexp.MustCompile(`\.(?P<ext>[A-Za-z0-9]+)$`)
}

func main() {
	uploadHandler := &UploadHandler{
		HostAddr:  "http://localhost:8081",
		UploadDir: "upload",
	}
	http.Handle("/upload", uploadHandler)

	browseHandler := &BrowseHandles{
		BrowseDir: "upload",
	}
	http.Handle("/", browseHandler)

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

	if r.Method == http.MethodPost {
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

		//fileLink := "http://" + h.HostAddr + "/" + header.Filename
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
	} else {
		http.Error(w, "Unsupported method", http.StatusBadRequest)
	}

}

func (b *BrowseHandles) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		path := b.BrowseDir
		if !strings.HasPrefix(path, "/") {
			currentDir, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			path = currentDir + "/" + path
		}

		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		res := []FileInfo{}
		for _, file := range files {
			extension := ""
			extensionMutch := FILE_EXT_REGEXP.FindStringSubmatch(file.Name())
			if len(extensionMutch) > 1 {
				extension = extensionMutch[1]
			}
			fileInfo := FileInfo{
				Filename:  file.Name(),
				Extension: extension,
				Size:      file.Size(),
			}
			res = append(res, fileInfo)
		}

		json.NewEncoder(w).Encode(res)

	} else {
		http.Error(w, "Unsupported method", http.StatusBadRequest)
	}

}
