package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Service struct {
	connectionString string
	students         IStore
	sync.RWMutex
}

func NewService(connectionString string, students IStore) *Service {
	return &Service{
		connectionString: connectionString,
		students:         students,
	}
}

func (s *Service) ListenAndServe() error {
	r := mux.NewRouter()

	r.HandleFunc("/student", logs(auth(s.PostStudent))).Methods("POST")
	r.HandleFunc("/student", logs(auth(s.GetStudents))).Methods("GET")
	r.HandleFunc("/student/{name}", logs(auth(s.GetStudent))).Methods("GET")
	r.HandleFunc("/student/{name}", logs(auth(s.PutStudent))).Methods("PUT")
	r.HandleFunc("/student/{name}", logs(auth(s.DeleteStudent))).Methods("DELETE")

	log.Printf("Starting server on %s", s.connectionString)
	err := http.ListenAndServe(s.connectionString, r)
	if err != nil {
		return err
	}
	return nil
}

func logs(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		log.Printf("%s %s", method, path)
		handlerFunc(w, r)
		return
	}
}

func auth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Please supply and Authorization token", http.StatusUnauthorized)
			return
		}
		handlerFunc(w, r)
		return
	}
}
