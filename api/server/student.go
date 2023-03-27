package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

type Student struct {
	Name        string
	Description string
	Subjects    []string
}

func (s *Service) GetStudents(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	err := json.NewEncoder(w).Encode(s.students)
	if err != nil {
		log.Println(err)
	}
}

func (s *Service) PostStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(student.Name)) {
		http.Error(w, "student names cannot contain whitespace", 400)
		return
	}

	s.Lock()
	defer s.Unlock()

	if s.studentExists(student.Name) {
		http.Error(w, fmt.Sprintf("student %s already exists", student.Name), http.StatusBadRequest)
		return
	}

	if err := s.students.put(student.Name, student); err != nil {
		log.Fatal("Failed storing ", student.Name)
	}
	log.Printf("added Student: %s", student.Name)
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		log.Printf("error sending response - %s", err)
	}
}

func (s *Service) PutStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentName := vars["name"]
	if studentName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	var student Student
	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	s.Lock()
	defer s.Unlock()

	if !s.studentExists(studentName) {
		log.Printf("student %s does not exist", studentName)
		http.Error(w, fmt.Sprintf("student %v does not exist", studentName), http.StatusBadRequest)
		return
	}

	if err := s.students.put(student.Name, student); err != nil {
		log.Fatal("Failed storing ", student.Name)
	}
	log.Printf("updated student: %s", student.Name)
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		log.Printf("error sending response - %s", err)
	}
}

func (s *Service) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentName := vars["name"]
	if studentName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	s.Lock()
	defer s.Unlock()

	if !s.studentExists(studentName) {
		http.Error(w, fmt.Sprintf("student %s does not exists", studentName), http.StatusNotFound)
		return
	}

	if err := s.students.delete(studentName); err != nil {
		log.Fatal("Unable to delete ", studentName)
	}

	_, err := fmt.Fprintf(w, "Deleted student with name %s", studentName)
	if err != nil {
		log.Println(err)
	}
}

func (s *Service) GetStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentName := vars["name"]
	if studentName == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	s.RLock()
	defer s.RUnlock()
	if !s.studentExists(studentName) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	student, _ := s.students.get(studentName)
	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Service) studentExists(studentName string) bool {
	if _, err := s.students.get(studentName); err == nil {
		return true
	}
	return false
}
