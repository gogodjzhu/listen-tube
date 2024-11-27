package auth

import (
	"errors"
	"time"

	"github.com/gogodjzhu/listen-tube/internal/pkg/db/dao"
	"golang.org/x/crypto/bcrypt"
)

// Authenticator authenticates the user
type Authenticator interface {
	// Authenticate authenticates the user
	Authenticate(username, password string) (*dao.User, error)
	// Register registers a new user
	Register(username, password string) error
	// ChangePassword changes the password of the user
	ChangePassword(username, oldPassword, newPassword string) error
}

// Authorize checks if the user is authorized to perform the action
type Authorizer interface {
	// Authorize checks if the user is authorized to perform the action
	Authorize(username, action string) (bool, error)
}

type AuthService struct {
	userMapper *dao.UserMapper
}

func NewAuthService(mapper *dao.UnionMapper) (*AuthService, error) {
	return &AuthService{
		userMapper: mapper.UserMapper,
	}, nil
}

// Authenticate authenticates the user
func (s *AuthService) Authenticate(username, password string) (*dao.User, error) {
	users, err := s.userMapper.Select(&dao.User{Name: username})
	if (err != nil) {
		return nil, err
	}
	if len(users) != 1 {
		return nil, errors.New("invalid user or password")
	}
	user := users[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.Credit), []byte(password))
	if err != nil {
		return nil, errors.New("invalid user or password")
	}
	return user, nil
}

// Register registers a new user
func (s *AuthService) Register(username, password string) error {
	users, err := s.userMapper.Select(&dao.User{Name: username})
	if err != nil {
		return err
	}
	if len(users) > 0 {
		return errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &dao.User{
		Name:     username,
		Credit:   string(hashedPassword),
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	_, err = s.userMapper.Insert(user)
	return err
}

// ChangePassword changes the password of the user
func (s *AuthService) ChangePassword(username, oldPassword, newPassword string) error {
	users, err := s.userMapper.Select(&dao.User{Name: username})
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("user not found")
	}
	user := users[0]
	err = bcrypt.CompareHashAndPassword([]byte(user.Credit), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Credit = string(hashedPassword)
	user.UpdateAt = time.Now()
	_, err = s.userMapper.Update(&dao.User{ID: user.ID}, user)
	return err
}