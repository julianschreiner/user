package user

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUserIfNotExisting(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, user *User) error
	GetActiveUserByEmail(email string) (*User, error)
	GetAllUsers(ctx context.Context, query *User) ([]*User, error)
	GetUserById(ctx context.Context, id uint64) (*User, error)
	SaveUser(ctx context.Context, user *User) error
	UpdateUserGroups(ctx context.Context, user *User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUserIfNotExisting(ctx context.Context, user *User) error {

	err := r.db.Where("email = ?", user.Email).
		Or("id = ?", user.ID).
		Attrs(user).
		FirstOrCreate(user).Error

	if err != nil {
		return errors.New("InternalServerError")
	}

	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, user *User) error {
	err := r.db.Delete(user).Error

	if err != nil {
		return errors.New("InternalServerError")
	}

	return nil
}

func (r *userRepo) GetActiveUserByEmail(email string) (*User, error) {
	u := User{
		Email:  email,
		Active: true,
	}

	err := r.db.Preload("Groups").Preload("Groups.Permissions").Where(&u).Take(&u).Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *userRepo) GetAllUsers(ctx context.Context, query *User) ([]*User, error) {
	var users []*User
	err := r.db.Where(query).Preload("Groups").Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("InternalServerError")
	}

	return users, nil
}

func (r *userRepo) GetUserById(ctx context.Context, id uint64) (*User, error) {
	var u User
	err := r.db.Preload("Groups").Preload("Groups.Permissions").First(&u, id).Error

	if err != nil {
		return nil, errors.New("NotFound")
	}

	return &u, nil
}

func (r *userRepo) SaveUser(ctx context.Context, user *User) error {
	err := r.db.Save(user).Error

	if err != nil {
		return errors.New("InternalServerError")
	}

	return nil
}

func (r *userRepo) UpdateUserGroups(ctx context.Context, user *User) error {
	/*err := r.db.Model(user).Association("Groups").Replace(user.Groups)

	if err != nil {
		return errors.New("InternalServerError")
	}*/

	return nil
}
