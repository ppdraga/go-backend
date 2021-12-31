package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"red/metrics"
)

var (
	measurable = metrics.MeasurableHandler
	router     = mux.NewRouter()
	web        = http.Server{
		Handler: router,
		//Addr:    ":8001",
	}
)

func init() {
	router.
		HandleFunc("/identity", measurable(GetIdentityHandler)).
		Methods(http.MethodGet)
}

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9091", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

func GetIdentityHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("token") == "admin_secret_token" {
		fmt.Println("auth ok")
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println("auth bad")
	w.WriteHeader(http.StatusUnauthorized)
}
