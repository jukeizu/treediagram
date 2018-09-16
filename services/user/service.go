package user

import (
	"context"

	pb "github.com/jukeizu/treediagram/api/user"
)

type service struct {
	UserStorage UserStorage
}

func NewService(userStorage UserStorage) pb.UserServer {
	return &service{userStorage}
}

func (s service) Preference(ctx context.Context, req *pb.PreferenceRequest) (*pb.PreferenceReply, error) {
	preference, err := s.UserStorage.Preference(req)

	return &pb.PreferenceReply{Preference: preference}, err
}

func (s service) SetServer(ctx context.Context, req *pb.SetServerRequest) (*pb.PreferenceReply, error) {
	preference, err := s.UserStorage.SetServer(req)

	return &pb.PreferenceReply{Preference: preference}, err
}
