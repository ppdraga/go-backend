package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadHandler(t *testing.T) {

	// Тест из методички падает с ошибкой:
	// panic: runtime error: invalid memory address or nil pointer dereference

	file, err := os.Open("testfile.txt")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok!")
	}))
	defer ts.Close()

	uploadHandler := &UploadHandler{
		HostAddr:  ts.URL,
		UploadDir: "upload",
	}
	uploadHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `testfile.txt`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
