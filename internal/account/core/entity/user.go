package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Username    string
	Password    string
	Email       string
	PhoneNumber string
	CreatedAt   time.Time
	updatedAt   time.Time
}
