package user

import (
	"gorm.io/gorm"
	pb "uas/user"
)

type User struct {
	gorm.Model
	Email    string
	Active   bool
	Role     uint64
	Forename string
	Surname  string
	Dob      string
}

func (u *User) ToPb() *pb.Client {
	user := &pb.Client{
		Id:        uint64(u.ID),
		Email:     u.Email,
		Active:    u.Active,
		Role:      u.Role,
		Forename:  u.Forename,
		Surname:   u.Surname,
		Dob:       u.Dob,
		CreatedAt: u.CreatedAt.String(),
	}

	return user
}
