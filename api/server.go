package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router
	list []Item
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
		list:   []Item{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/list", s.addItem()).Methods("POST")
	s.HandleFunc("/list", s.listItems()).Methods("GET")
	s.HandleFunc("/list/{id}", s.removeItem()).Methods("DELETE")
}

func (s *Server) addItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i Item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		s.list = append(s.list, i)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.list); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) removeItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr, _ := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		for i, item := range s.list {
			if item.ID == id {
				s.list = append(s.list[:i], s.list[i+1:]...)
				break
			}
		}
	}
}
