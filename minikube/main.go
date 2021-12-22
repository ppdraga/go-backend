package main

import (
	"context"
	"encoding/json"
	"fmt"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	_ "google.golang.org/grpc/reflection"
	"log"
	proto "minikube/pkg/addservice"
	"minikube/version"
	"net"
	"net/http"
	"sync"
)

type server struct {
	proto.UnimplementedAddServiceServer
}

func main() {
	wg := sync.WaitGroup{}
	fmt.Println("Ver: ", version.Version)
	fmt.Println("Commit: ", version.Commit)
	// Web server
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, _ = fmt.Fprintf(w, "Simple go web app!!!\n")
	})
	http.HandleFunc("/__heartbeat__", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "heartbeat endpoint\n")
	})
	http.HandleFunc("/__version__", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"version": version.Version,
			"commit":  version.Commit,
			"build":   version.Build,
		})
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := http.ListenAndServe(":8080", nil); e != nil {
			log.Fatal(e)
		}
	}()
	log.Println("http server up and running on port 8080")

	// GRPC server
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	proto.RegisterAddServiceServer(srv, &server{})
	reflection.Register(srv)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := srv.Serve(listener); e != nil {
			log.Fatal(e)
		}
	}()
	log.Println("grpc server up and running on port 8082")

	wg.Wait()
}

func (s *server) Add(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	a, b := request.GetA(), request.GetB()
	result := a + b
	return &proto.Response{Result: result}, nil
}

func (s *server) Multiply(ctx context.Context, request *proto.Request) (*proto.Response, error) {
	a, b := request.GetA(), request.GetB()
	result := a * b
	return &proto.Response{Result: result}, nil
}
