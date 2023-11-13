package server

import (
	"context"
	"net/http"
	"time"
	"usernet/internal/app/repos/aroundment"

	"usernet/internal/app/repos/user"
	"usernet/internal/app/starter"
)

var _ starter.APIServer = &Server{}

type Server struct {
	srv http.Server
	us  *user.Users
	as  *aroundment.Aroundments
}

func NewServer(addr string, h http.Handler) *Server {
	s := &Server{}

	s.srv = http.Server{
		Addr:              addr,
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
	}
	return s
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start(us *user.Users, as *aroundment.Aroundments) {
	s.us = us
	s.as = as
	go s.srv.ListenAndServe()
}
