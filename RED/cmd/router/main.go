package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"red/metrics"
	"time"
)

var (
	db         *sql.DB
	measurable = metrics.MeasurableHandler
	router     = mux.NewRouter()
	web        = http.Server{
		Handler: router,
		//Addr:    ":8002",
	}
)

func init() {
	router.
		HandleFunc("/entities", measurable(ListEntitiesHandler)).
		Methods(http.MethodGet)
	router.
		HandleFunc("/entity", measurable(AddEntityHandler)).
		Methods(http.MethodPost)
	var err error
	db, err = sql.Open("mysql", "root:qwerty@tcp(mysql:3306)/test")
	//db, err = sql.Open("mysql", "root:qwerty@tcp(localhost:3306)/test")
	if err != nil {
		panic(err)
	}
}

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != http.ErrServerClosed {
			panic(fmt.Errorf("error on listen and serve: %v", err))
		}
	}()
	if err := web.ListenAndServe(); err != http.ErrServerClosed {
		panic(fmt.Errorf("error on listen and serve: %v", err))
	}
}

const sqlInsertEntity = "INSERT INTO entities(id, data) VALUES (?, ?)"

func AddEntityHandler(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	id := r.FormValue("id")
	data := r.FormValue("data")
	fmt.Println("token", token)
	fmt.Println("id", id)
	fmt.Println("data", data)

	res, err := http.Get(fmt.Sprintf("http://acl:80/identity?token=%s",
		//res, err := http.Get(fmt.Sprintf("http://localhost:8001/identity?token=%s",
		r.FormValue("token")))
	switch {
	case err != nil:
		fmt.Println("err", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	case res.StatusCode != http.StatusOK:
		fmt.Println("err", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	res.Body.Close()

	t := time.Now()
	q := "sqlInsertEntity"
	metrics.RequestsTotalSql.WithLabelValues(q).Inc()
	_, err = db.Exec(sqlInsertEntity, r.FormValue("id"), r.FormValue("data"))
	metrics.DurationSql.WithLabelValues(q).Observe(time.Since(t).Seconds())
	if err != nil {
		metrics.ErrorsTotalSql.WithLabelValues(q).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

const sqlSelectEntities = "SELECT id, data FROM entities"

type ListEntityItemResponse struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

func ListEntitiesHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	q := "sqlSelectEntities"
	metrics.RequestsTotalSql.WithLabelValues(q).Inc()
	rr, err := db.Query(sqlSelectEntities)
	metrics.DurationSql.WithLabelValues(q).Observe(time.Since(t).Seconds())
	if err != nil {
		metrics.ErrorsTotalSql.WithLabelValues(q).Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rr.Close()
	var ii []*ListEntityItemResponse
	for rr.Next() {
		i := &ListEntityItemResponse{}
		err = rr.Scan(&i.Id, &i.Data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ii = append(ii, i)
	}
	bb, err := json.Marshal(ii)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(bb)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
