package entity

import (
	"errors"
	"regexp"
	"strings"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewUser(name, email string) (*User, error) {
	user := &User{
		Name:  strings.TrimSpace(name),
		Email: strings.TrimSpace(strings.ToLower(email)),
	}
	
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	return user, nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	
	if u.Email == "" {
		return errors.New("email is required")
	}
	
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}
	
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
