package handlers

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	pb "uas/user"
	"uas/user/user"
)

const NormalUserId = 3

var logger2 log.Logger

// NewService returns a naive, stateless implementation of Service.
func NewService() pb.UserServer {
	ioWriter := log.New(os.Stdout, "\r\n", log.LstdFlags)
	ioWriter.Println("started user-uas-svc")

	dbHost := os.Getenv("DB_HOST")
	dbUserName := os.Getenv("DB_USER")
	dbSecret := os.Getenv("DB_SECRET")
	dbName := os.Getenv("DB_NAME")

	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=UTC", dbUserName, dbSecret, dbHost, dbName)), &gorm.Config{
		Logger: logger.New(ioWriter,
			logger.Config{
				SlowThreshold: time.Millisecond * 200,
				LogLevel:      0,
			},
		),
	})
	if err != nil {
		println("database is not reachable", "error", err)
		os.Exit(3)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		println("failed to migrate db", "error", err)
		os.Exit(7)
	}

	/* DOMAIN LOGIC */
	userRepository := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepository, &logger2)

	return userService{
		userManger: userSvc,
	}
}

type userService struct {
	userManger user.UserService
}

func (s userService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	input := user.User{
		Email:    in.Email,
		Active:   true,
		Role:     in.Role,
		Forename: in.Forename,
		Surname:  in.Surname,
		Dob:      in.Dob,
	}

	resp, err := s.userManger.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		User: resp.ToPb(),
	}, nil
}

func (s userService) GetAllUserInformation(ctx context.Context, in *pb.GetAllUserInformationRequest) (*pb.GetAllUserInformationResponse, error) {
	resp, err := s.userManger.GetAllUsersByRole(ctx, NormalUserId)
	if err != nil {
		return nil, err
	}

	clients := make([]*pb.Client, len(resp))
	for i, val := range resp {
		clients[i] = val.ToPb()
	}

	return &pb.GetAllUserInformationResponse{
		Users: clients,
	}, nil
}

func (s userService) GetUserInformation(ctx context.Context, in *pb.GetUserInformationRequest) (*pb.GetUserInformationResponse, error) {
	resp, err := s.userManger.GetUserById(ctx, uint(in.Id))
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return &pb.GetUserInformationResponse{
			User: nil,
		}, nil
	}

	return &pb.GetUserInformationResponse{
		User: resp.ToPb(),
	}, nil
}

func (s userService) GetUserInformationEmail(ctx context.Context, in *pb.GetUserInformationEmailRequest) (*pb.GetUserInformationEmailResponse, error) {
	println("received request getUserInformation Email")
	println("input: " + in.String())

	resp, err := s.userManger.GetActiveUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return &pb.GetUserInformationEmailResponse{
			User: nil,
		}, nil
	}

	return &pb.GetUserInformationEmailResponse{
		User: resp.ToPb(),
	}, nil
}

func (s userService) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.userManger.DeleteUserById(ctx, uint(in.Id))
	success := true
	if err != nil {
		success = false
	}

	return &pb.DeleteUserResponse{Success: success}, nil
}
