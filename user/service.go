package user

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
)

type service struct {
	Repository Repository
}

func NewService(repository Repository) pb.UserServer {
	return &service{repository}
}

func (s service) Preference(ctx context.Context, req *pb.PreferenceRequest) (*pb.PreferenceReply, error) {
	preference, err := s.Repository.Preference(req)

	return &pb.PreferenceReply{Preference: preference}, err
}

func (s service) SetServer(ctx context.Context, req *pb.SetServerRequest) (*pb.PreferenceReply, error) {
	preference, err := s.Repository.SetServer(req)

	return &pb.PreferenceReply{Preference: preference}, err
}
