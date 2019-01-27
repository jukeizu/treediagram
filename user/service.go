package user

import (
	"context"

	"github.com/jukeizu/treediagram/api/protobuf-spec/userpb"
)

type service struct {
	Repository Repository
}

func NewService(repository Repository) userpb.UserServer {
	return &service{repository}
}

func (s service) Preference(ctx context.Context, req *userpb.PreferenceRequest) (*userpb.PreferenceReply, error) {
	preference, err := s.Repository.Preference(req)

	return &userpb.PreferenceReply{Preference: preference}, err
}

func (s service) SetServer(ctx context.Context, req *userpb.SetServerRequest) (*userpb.PreferenceReply, error) {
	preference, err := s.Repository.SetServer(req)

	return &userpb.PreferenceReply{Preference: preference}, err
}
