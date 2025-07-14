package domain

import (
	"errors"
	"fmt"
)

type User struct {
	BaseModel
	Email     string  `gorm:"type:varchar(255);unique;not null"`
	FirstName string  `gorm:"type:varchar(32);not null"`
	LastName  string  `gorm:"type:varchar(32);not null"`
	Tasks     []*Task `gorm:"foreignKey:AssigneeID"`
}

var ErrValidation = errors.New("validation error")

func NewUser(userId, email, firstName, lastName string) (*User, error) {
	if userId == "" {
		return nil, fmt.Errorf("%w: userId is required", ErrValidation)
	}
	if email == "" {
		return nil, fmt.Errorf("%w: email is required", ErrValidation)
	}
	if firstName == "" {
		return nil, fmt.Errorf("%w: firstName is required", ErrValidation)
	}
	if lastName == "" {
		return nil, fmt.Errorf("%w: lastName is required", ErrValidation)
	}

	return &User{
		BaseModel: BaseModel{
			ID: userId,
		},
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
