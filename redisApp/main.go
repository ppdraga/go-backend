package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	Name   string
	Email  string
	Pass   string
	Active bool
}

var users map[uuid.UUID]*User
var redisClient *RedisClient

func init() {
	users = make(map[uuid.UUID]*User)
	rand.Seed(time.Now().UnixNano())

	const (
		redisHost = "localhost"
		redisPort = "6379"
	)
	var err error
	redisClient, err = NewRedisClient(redisHost, redisPort)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	defer redisClient.Close()

	router := mux.NewRouter()
	router.HandleFunc("/singup", SignUpEndpoint).Methods("POST")
	router.HandleFunc("/check", CheckEndpoint).Methods("GET")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Print("Server Started")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server Failed: %+v", err)
	}

}

type SignUpMsg struct {
	Name  string `json:"name"`
	Email string `json:"e-mail"`
	Pass  string `json:"pass"`
}

func SignUpEndpoint(w http.ResponseWriter, r *http.Request) {
	var signUpMsg SignUpMsg
	err := json.NewDecoder(r.Body).Decode(&signUpMsg)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "internal error: %+v\n", err)
		return
	}
	log.Println(signUpMsg)
	id := uuid.New()
	users[id] = &User{
		Name:   signUpMsg.Name,
		Email:  signUpMsg.Email,
		Pass:   signUpMsg.Pass,
		Active: false,
	}
	checkMsg := RandStringBytes(16)
	log.Println(users)
	err = redisClient.Set(context.Background(), checkMsg, id.String(), time.Hour).Err()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "internal error: %+v\n", err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "Please, check e-mail %s and confirm your account within one hour!\nLink http://localhost:8080/check?msg=%s\n", signUpMsg.Email, checkMsg)
}

func CheckEndpoint(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		w.WriteHeader(500)
		fmt.Fprintf(w, "internal error: Couldn't parse msg param\n")
		return
	}

	userIdStr, err := redisClient.GetRecord(msg)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "internal error: %+v\n", err)
		return
	}
	userId, _ := uuid.Parse(string(userIdStr))
	user, ok := users[userId]
	if !ok {
		w.WriteHeader(500)
		w.Write([]byte("Check failed!\n"))
		return
	}
	user.Active = true
	w.WriteHeader(200)
	w.Write([]byte("Check OK!\n"))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
