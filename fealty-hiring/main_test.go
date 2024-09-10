package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudent).Methods("GET")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", generateSummary).Methods("GET")
	return r
}

func TestCreateStudent(t *testing.T) {
	store = NewStudentStore()
	router := setupRouter()

	student := Student{Name: "John Doe", Age: 20, Email: "john@example.com"}
	body, _ := json.Marshal(student)

	req, _ := http.NewRequest("POST", "/students", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response map[string]int
	json.Unmarshal(rr.Body.Bytes(), &response)

	if id, exists := response["id"]; !exists || id != 1 {
		t.Errorf("expected id 1, got %v", id)
	}
}

func TestGetStudents(t *testing.T) {
	store = NewStudentStore()
	store.Add(Student{Name: "John Doe", Age: 20, Email: "john@example.com"})
	store.Add(Student{Name: "Jane Doe", Age: 22, Email: "jane@example.com"})

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/students", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var students []Student
	json.Unmarshal(rr.Body.Bytes(), &students)

	if len(students) != 2 {
		t.Errorf("expected 2 students, got %d", len(students))
	}
}

func TestGetStudent(t *testing.T) {
	store = NewStudentStore()
	store.Add(Student{Name: "John Doe", Age: 20, Email: "john@example.com"})

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/students/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var student Student
	json.Unmarshal(rr.Body.Bytes(), &student)

	if student.Name != "John Doe" {
		t.Errorf("expected student name John Doe, got %v", student.Name)
	}
}

func TestUpdateStudent(t *testing.T) {
	store = NewStudentStore()
	store.Add(Student{Name: "John Doe", Age: 20, Email: "john@example.com"})

	router := setupRouter()

	updatedStudent := Student{Name: "John Updated", Age: 21, Email: "john.updated@example.com"}
	body, _ := json.Marshal(updatedStudent)

	req, _ := http.NewRequest("PUT", "/students/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	student, _ := store.Get(1)
	if student.Name != "John Updated" || student.Age != 21 ||
		student.Email != "john.updated@example.com" {
		t.Errorf("student was not updated correctly")
	}
}

func TestDeleteStudent(t *testing.T) {
	store = NewStudentStore()
	store.Add(Student{Name: "John Doe", Age: 20, Email: "john@example.com"})

	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/students/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	_, exists := store.Get(1)
	if exists {
		t.Errorf("student was not deleted")
	}
}

func TestGenerateSummary(t *testing.T) {
	store = NewStudentStore()
	store.Add(Student{Name: "John Doe", Age: 20, Email: "john@example.com"})

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/students/1/summary", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	json.Unmarshal(rr.Body.Bytes(), &response)

	expectedSummary := "Student John Doe is 20 years old and can be contacted at john@example.com."
	if summary, exists := response["summary"]; !exists || summary != expectedSummary {
		t.Errorf("expected summary %s, got %s", expectedSummary, summary)
	}
}
