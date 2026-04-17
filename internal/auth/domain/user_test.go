package domain

import (
	"errors"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		userName string
		wantErr  error
	}{
		{"valid user", "admin@test.com", "password123", "Admin", nil},
		{"empty email", "", "password123", "Admin", ErrEmptyEmail},
		{"empty name", "admin@test.com", "password123", "", ErrEmptyName},
		{"empty password", "admin@test.com", "", "Admin", ErrEmptyPassword},
		{"short password", "admin@test.com", "12345", "Admin", ErrPasswordTooShort},
		{"password exactly 6 chars", "admin@test.com", "123456", "Admin", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.userName)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewUser() error = %v, want %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && user == nil {
				t.Fatal("expected user to be non-nil")
			}

			if tt.wantErr == nil && user != nil {
				if user.Email != tt.email {
					t.Errorf("Email = %q, want %q", user.Email, tt.email)
				}
				if user.Name != tt.userName {
					t.Errorf("Name = %q, want %q", user.Name, tt.userName)
				}
				if user.Password == tt.password {
					t.Error("Password should be hashed, not plain text")
				}
				if user.Password == "" {
					t.Error("Password hash should not be empty")
				}
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	user, err := NewUser("admin@test.com", "password123", "Admin")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"correct password", "password123", true},
		{"wrong password", "wrongpassword", false},
		{"empty password", "", false},
		{"similar password", "password124", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := user.CheckPassword(tt.password); got != tt.want {
				t.Errorf("CheckPassword(%q) = %v, want %v", tt.password, got, tt.want)
			}
		})
	}
}

func TestNewUser_PasswordHashing(t *testing.T) {
	user1, _ := NewUser("a@test.com", "samepassword", "A")
	user2, _ := NewUser("b@test.com", "samepassword", "B")

	if user1.Password == user2.Password {
		t.Error("same password should produce different hashes (bcrypt uses random salt)")
	}

	if !user1.CheckPassword("samepassword") {
		t.Error("user1 should validate with original password")
	}
	if !user2.CheckPassword("samepassword") {
		t.Error("user2 should validate with original password")
	}
}
