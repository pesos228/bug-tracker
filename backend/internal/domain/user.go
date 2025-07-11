package domain

import "errors"

type User struct {
	BaseModel
	Email     string `gorm:"type:varchar(255);unique;not null"`
	FirstName string `gorm:"type:varchar(32);not null"`
	LastName  string `gorm:"type:varchar(32);not null"`
	Tasks     []*Task
}

var ErrValidation = errors.New("validation error")

func NewUser(userId, email, firstName, lastName string) (*User, error) {
	if userId == "" || email == "" || firstName == "" || lastName == "" {
		return nil, ErrValidation
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
