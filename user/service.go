package user

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/protobuf-spec/user"
)

type service struct {
	UserDb UserDb
}

func NewService(userDb UserDb) pb.UserServer {
	return &service{userDb}
}

func (s service) Preference(ctx context.Context, req *pb.PreferenceRequest) (*pb.PreferenceReply, error) {
	preference, err := s.UserDb.Preference(req)

	return &pb.PreferenceReply{Preference: preference}, err
}

func (s service) SetServer(ctx context.Context, req *pb.SetServerRequest) (*pb.PreferenceReply, error) {
	preference, err := s.UserDb.SetServer(req)

	return &pb.PreferenceReply{Preference: preference}, err
}
