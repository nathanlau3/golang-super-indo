package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrEmailAlreadyExists = errors.New("email sudah terdaftar")
	ErrEmptyEmail         = errors.New("email wajib diisi")
	ErrEmptyPassword      = errors.New("password wajib diisi")
	ErrPasswordTooShort   = errors.New("password minimal 6 karakter")
	ErrEmptyName          = errors.New("nama wajib diisi")
)

type User struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(email, password, name string) (*User, error) {
	if email == "" {
		return nil, ErrEmptyEmail
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	if password == "" {
		return nil, ErrEmptyPassword
	}
	if len(password) < 6 {
		return nil, ErrPasswordTooShort
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:    email,
		Password: string(hashed),
		Name:     name,
	}, nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}
