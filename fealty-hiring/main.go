package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "hello"})
}

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type StudentStore struct {
	students map[int]Student
	mutex    sync.RWMutex
	nextID   int
}

func NewStudentStore() *StudentStore {
	return &StudentStore{
		students: make(map[int]Student),
		nextID:   1,
	}
}

func (s *StudentStore) Add(student Student) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	student.ID = s.nextID
	s.students[student.ID] = student
	s.nextID++
	return student.ID
}

func (s *StudentStore) Get(id int) (Student, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	student, ok := s.students[id]
	return student, ok
}

func (s *StudentStore) GetAll() []Student {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	students := make([]Student, 0, len(s.students))
	for _, student := range students {
		students = append(students, student)
	}
	return students
}

func (s *StudentStore) Update(id int, student Student) bool {

	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.students[id]; !ok {
		return false
	}
	student.ID = id
	s.students[id] = student
	return true
}

func (s *StudentStore) Delete(id int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.students[id]; !ok {
		return false
	}
	delete(s.students, id)
	return true
}

var store = NewStudentStore()

func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student

	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	if student.Name == "" || student.Age <= 0 || student.Email == "" {
		http.Error(w, "Invalid Student Data", http.StatusBadRequest)
		return
	}

	id := store.Add(student)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	students := store.GetAll()
	json.NewEncoder(w).Encode(students)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student, ok := store.Get(id)
	if !ok {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid student id", http.StatusBadRequest)
		return
	}

	var student Student
	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if student.Name == "" || student.Age <= 0 || student.Email == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if !store.Update(id, student) {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	if !store.Delete(id) {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Student ID", http.StatusBadRequest)
		return
	}

	student, ok := store.Get(id)
	if !ok {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	summary := fmt.Sprintf(
		"Student %s is %d years old and can be contacted at %s.",
		student.Name,
		student.Age,
		student.Email,
	)
	json.NewEncoder(w).Encode(map[string]string{"summary": summary})
}

func main() {
	seedFlag := flag.Bool("seed", false, "Seed the database with sample data")
	flag.Parse()

	if *seedFlag {
		SeedData()
		return
	}
	r := mux.NewRouter()

	r.HandleFunc("/healthcheck", healthCheck).Methods("GET")
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudent).Methods("GET")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", generateSummary).Methods("GET")

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
