package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"usernet/internal/app/repos/aroundment"
	"usernet/internal/db/mem/aroundmentmemstore"

	"usernet/internal/app/repos/user"
	"usernet/internal/db/mem/usermemstore"
)

func TestRouter_CreateUser(t *testing.T) {

	ust := usermemstore.NewUsers()
	ast := aroundmentmemstore.NewAroundments()

	us := user.NewUsers(ust)
	as := aroundment.NewAroundments(ast)

	rt := NewRouter(us, as)

	hts := httptest.NewServer(rt)

	r, _ := http.NewRequest("POST", hts.URL+"/user/create", strings.NewReader(`{"name":"user123"}`))
	r.SetBasicAuth("admin", "admin")

	cli := hts.Client()

	resp, err := cli.Do(r)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Error("status wrong:", resp.StatusCode)
	}

	// (&http.Client{}).Get(httptest.NewServer(nil).URL)
}

func TestRouter_CreateUser2(t *testing.T) {
	ust := usermemstore.NewUsers()
	ast := aroundmentmemstore.NewAroundments()

	us := user.NewUsers(ust)
	as := aroundment.NewAroundments(ast)

	rt := NewRouter(us, as)

	h := rt.AuthMiddleware(http.HandlerFunc(rt.CreateUser)).ServeHTTP

	w := &httptest.ResponseRecorder{}
	r := httptest.NewRequest("POST", "/user/create", strings.NewReader(`{"name":"user123"}`))
	r.SetBasicAuth("admin", "admin")

	h(w, r)

	if w.Code != http.StatusCreated {
		t.Error("status wrong:", w.Code)
	}
}
