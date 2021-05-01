package handlers

import (
	"context"

	pb "uas/user"
)

// NewService returns a naive, stateless implementation of Service.
func NewService() pb.UserServer {
	return userService{}
}

type userService struct{}

func (s userService) Information(ctx context.Context, in *pb.InformationRequest) (*pb.InformationResponse, error) {
	var resp pb.InformationResponse
	return &resp, nil
}
