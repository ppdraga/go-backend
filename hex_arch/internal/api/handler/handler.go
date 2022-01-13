package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"usernet/internal/app/repos/aroundment"

	"github.com/google/uuid"
	"usernet/internal/app/repos/user"
)

type Router struct {
	*http.ServeMux
	us *user.Users
	as *aroundment.Aroundments
}

func NewRouter(us *user.Users, as *aroundment.Aroundments) *Router {
	r := &Router{
		ServeMux: http.NewServeMux(),
		us:       us,
		as:       as,
	}
	r.Handle("/user/create", r.AuthMiddleware(http.HandlerFunc(r.CreateUser)))
	r.Handle("/user/read", r.AuthMiddleware(http.HandlerFunc(r.ReadUser)))
	r.Handle("/user/delete", r.AuthMiddleware(http.HandlerFunc(r.DeleteUser)))
	r.Handle("/user/search", r.AuthMiddleware(http.HandlerFunc(r.SearchUser)))

	r.Handle("/aroundment/create", r.AuthMiddleware(http.HandlerFunc(r.CreateAroundment)))
	r.Handle("/aroundment/read", r.AuthMiddleware(http.HandlerFunc(r.ReadAroundment)))
	r.Handle("/aroundment/delete", r.AuthMiddleware(http.HandlerFunc(r.DeleteAroundment)))
	r.Handle("/aroundment/search", r.AuthMiddleware(http.HandlerFunc(r.SearchAroundment)))

	r.Handle("/aroundment/add_user", r.AuthMiddleware(http.HandlerFunc(r.AddUserToAroundment)))
	r.Handle("/aroundment/delete_user", r.AuthMiddleware(http.HandlerFunc(r.DeleteUserFromAroundment)))
	r.Handle("/aroundment/list_users", r.AuthMiddleware(http.HandlerFunc(r.ListUsersInAroundment)))
	r.Handle("/aroundment/search_by_membership", r.AuthMiddleware(http.HandlerFunc(r.SearchAroundmentsByUserMembership)))

	return r
}

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Data       string    `json:"data"`
	Permission int       `json:"perms"`
}

type Aroundment struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

func (rt *Router) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if u, p, ok := r.BasicAuth(); !ok || !(u == "admin" && p == "admin") {
				http.Error(w, "unautorized", http.StatusUnauthorized)
				return
			}
			// r = r.WithContext(context.WithValue(r.Context(), 1, 0))
			next.ServeHTTP(w, r)
		},
	)
}

func (rt *Router) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	u := User{}
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	bu := user.User{
		Name: u.Name,
		Data: u.Data,
	}

	nbu, err := rt.us.Create(r.Context(), bu)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(
		User{
			ID:         nbu.ID,
			Name:       nbu.Name,
			Data:       nbu.Data,
			Permission: nbu.Permissions,
		},
	)
}

// user/read?uid=...
func (rt *Router) ReadUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	suid := r.URL.Query().Get("uid")
	if suid == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(suid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (uid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	nbu, err := rt.us.Read(r.Context(), uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading", http.StatusInternalServerError)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(
		User{
			ID:         nbu.ID,
			Name:       nbu.Name,
			Data:       nbu.Data,
			Permission: nbu.Permissions,
		},
	)
}

func (rt *Router) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	suid := r.URL.Query().Get("uid")
	if suid == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(suid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (uid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = rt.us.Delete(r.Context(), uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "ok")
}

// /user/search?q=...
func (rt *Router) SearchUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	users, err := rt.us.SearchUsers(r.Context(), q)
	if err != nil {
		http.Error(w, "error when reading", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(users)

}

func (rt *Router) CreateAroundment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	a := Aroundment{}
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	ba := aroundment.Aroundment{
		Name: a.Name,
		Type: aroundment.AroundmentType(a.Type),
	}

	nba, err := rt.as.Create(r.Context(), ba)
	if err != nil {
		http.Error(w, "error when creating", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(
		Aroundment{
			ID:   nba.ID,
			Name: nba.Name,
			Type: string(nba.Type),
		},
	)
}

// read?uid=...
func (rt *Router) ReadAroundment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	said := r.URL.Query().Get("aid")
	if said == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	aid, err := uuid.Parse(said)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (aid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	nba, err := rt.as.Read(r.Context(), aid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading", http.StatusInternalServerError)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(
		Aroundment{
			ID:   nba.ID,
			Name: nba.Name,
			Type: string(nba.Type),
		},
	)
}

func (rt *Router) DeleteAroundment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	said := r.URL.Query().Get("aid")
	if said == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	aid, err := uuid.Parse(said)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if (aid == uuid.UUID{}) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = rt.as.Delete(r.Context(), aid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
		} else {
			http.Error(w, "error when reading", http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "ok")
}

// /search?q=...
func (rt *Router) SearchAroundment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("q")
	if q == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	aroundments, err := rt.as.SearchAroundments(r.Context(), q)
	if err != nil {
		http.Error(w, "error when reading", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(aroundments)

}

func (rt *Router) AddUserToAroundment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (rt *Router) DeleteUserFromAroundment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (rt *Router) ListUsersInAroundment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (rt *Router) SearchAroundmentsByUserMembership(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}
