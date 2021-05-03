package user

import (
	"context"
	"github.com/pkg/errors"
	"log"
)

type UserService interface {
	GetAllUsersByRole(ctx context.Context, roleId uint64) ([]*User, error)
	CreateUser(ctx context.Context, user User) (*User, error)
	SaveUser(ctx context.Context, user *User) error
	DeleteUserById(ctx context.Context, userId uint) error
	GetUserById(ctx context.Context, userId uint) (*User, error)
	GetActiveUserByEmail(ctx context.Context, email string) (*User, error)
}

type userService struct {
	userRepository UserRepository
	logger         *log.Logger
}

func NewUserService(userRepository UserRepository, log *log.Logger) UserService {
	return &userService{
		userRepository: userRepository,
		logger:         log,
	}
}

func (s *userService) CreateUser(ctx context.Context, user User) (*User, error) {
	err := s.userRepository.CreateUserIfNotExisting(ctx, &user)
	if err != nil {
		return nil, errors.New("InternalServerError")
	}

	return &user, nil
}

func (s *userService) DeleteUserById(ctx context.Context, userId uint) error {
	u, err := s.userRepository.GetUserById(ctx, uint64(userId))
	if err != nil {
		return err
	}

	/*check permission*/

	return s.userRepository.DeleteUser(ctx, u)
}

func (s *userService) GetActiveUserByEmail(ctx context.Context, email string) (*User, error) {
	// TODO CheckPermission out of abundance of precaution
	return s.userRepository.GetActiveUserByEmail(email)
}

func (s *userService) GetAllUsersByRole(ctx context.Context, roleId uint64) ([]*User, error) {
	return s.userRepository.GetAllUsers(ctx, &User{Role: roleId})
}

func (s *userService) GetUserById(ctx context.Context, userId uint) (*User, error) {
	// TODO CheckPermission out of abundance of precaution
	return s.userRepository.GetUserById(ctx, uint64(userId))
}

func (s *userService) SaveUser(ctx context.Context, user *User) error {
	// TODO CheckPermission out of abundance of precaution
	return s.userRepository.SaveUser(ctx, user)
}
