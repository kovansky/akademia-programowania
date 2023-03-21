package main

import (
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"github.com/grupawp/appdispatcher"
	"log"
)

func main() {
	code, err := appdispatcher.Submit(NewStudent("Stanisław", "Kowański"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(code)
}

type Student struct {
	applicationID uuid.UUID

	FirstName string
	LastName  string
}

func NewStudent(firstName, lastName string) Student {
	return Student{
		FirstName:     firstName,
		LastName:      lastName,
		applicationID: uuid.New(),
	}
}

func (s Student) ApplicationID() string {
	return s.applicationID.String()
}

func (s Student) FullName() string {
	return s.FirstName + " " + s.LastName
}
