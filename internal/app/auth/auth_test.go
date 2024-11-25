package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db"
	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
)

var fixedTime = time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)

func setupSuite(t *testing.T) func(t *testing.T) {
	// Return a function to teardown the test suite
	return func(t *testing.T) {
		// Teardown: delete test db files
		files, err := filepath.Glob("/tmp/listen-tube-unit-test-*.db")
		if err != nil {
			t.Fatalf("Failed to list test db files: %v", err)
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				t.Fatalf("Failed to delete test db file %s: %v", file, err)
			}
		}
	}
}

func setupTest(t *testing.T, s *AuthService) func(t *testing.T) {
	// Insert initial entities
	user := &dao.User{
		Credit:   "$2a$10$P/2om2cYvUSHBBwmjZqm5OVc4vR96zYmz5bxYf5u6126ziHVD7py.", // bcrypt hash for "password"
		Name:     "TestUser",
		CreateAt: fixedTime,
		UpdateAt: fixedTime,
	}
	_, err := s.userMapper.Insert(user)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	// Return a function to teardown the test
	return func(t *testing.T) {
		// Teardown: clean up the database
		s.userMapper.DB.Exec("DELETE FROM " + user.TableName())
	}
}

func MockAuthService() *AuthService {
	conf := &db.Config{
		DSN: fmt.Sprintf("/tmp/listen-tube-unit-test-%d.db", time.Now().UnixNano()),
	}
	ds, err := db.NewDatabaseSource(conf)
	if err != nil {
		panic(err)
	}

	unionMapper, err := dao.NewUnionMapper(ds)
	if err != nil {
		panic(err)
	}

	authService, err := NewAuthService(unionMapper)
	if err != nil {
		panic(err)
	}

	return authService
}

func TestAuthService_Authenticate(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name     string
		username string
		password string
		want     bool
		wantErr  bool
	}{
		{
			name:     "Valid credentials",
			username: "TestUser",
			password: "password",
			want:     true,
			wantErr:  false,
		},
		{
			name:     "Invalid password",
			username: "TestUser",
			password: "wrongpassword",
			want:     false,
			wantErr:  true,
		},
		{
			name:     "Non-existent user",
			username: "NonExistentUser",
			password: "password",
			want:     false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockAuthService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			got, err := s.Authenticate(tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AuthService.Authenticate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "Register new user",
			username: "NewUser",
			password: "newpassword",
			wantErr:  false,
		},
		{
			name:     "Register existing user",
			username: "TestUser",
			password: "password",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockAuthService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			if err := s.Register(tt.username, tt.password); (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_ChangePassword(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	tests := []struct {
		name        string
		username    string
		oldPassword string
		newPassword string
		wantErr     bool
	}{
		{
			name:        "Change password with valid old password",
			username:    "TestUser",
			oldPassword: "password",
			newPassword: "newpassword",
			wantErr:     false,
		},
		{
			name:        "Change password with invalid old password",
			username:    "TestUser",
			oldPassword: "wrongpassword",
			newPassword: "newpassword",
			wantErr:     true,
		},
		{
			name:        "Change password for non-existent user",
			username:    "NonExistentUser",
			oldPassword: "password",
			newPassword: "newpassword",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := MockAuthService()
			teardownTest := setupTest(t, s)
			defer teardownTest(t)

			if err := s.ChangePassword(tt.username, tt.oldPassword, tt.newPassword); (err != nil) != tt.wantErr {
				t.Errorf("AuthService.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}