package main

import (
	"fmt"
	"math/rand"
	"time"
)

var firstNames = []string{
	"John",
	"Jane",
	"Michael",
	"Emily",
	"David",
	"Sarah",
	"Robert",
	"Emma",
	"William",
	"Olivia",
}

var lastNames = []string{
	"Smith",
	"Johnson",
	"Brown",
	"Taylor",
	"Miller",
	"Anderson",
	"Wilson",
	"Moore",
	"Jackson",
	"Martin",
}

func SeedData() {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		student := generateRandomStudent()
		addStudent(student)
	}

	fmt.Println("Seeding completed successfully!")
}

func generateRandomStudent() Student {
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	name := firstName + " " + lastName
	age := rand.Intn(10) + 18 // Random age between 18 and 27
	email := fmt.Sprintf("%s.%s@example.com", firstName, lastName)

	return Student{
		Name:  name,
		Age:   age,
		Email: email,
	}
}

func addStudent(student Student) {
	id := store.Add(student)
	fmt.Printf("Added student: %s (ID: %d)\n", student.Name, id)
}
