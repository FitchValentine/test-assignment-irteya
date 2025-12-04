package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Firstname string
	Lastname  string
	Fullname  string
	Age       int
	IsMarried bool
	Password  string
}

