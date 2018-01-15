package main

type Service interface {
	Request(TreediagramRequest) (string, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Request(treediagramRequest TreediagramRequest) (string, error) {
	return treediagramRequest.Content, nil
}
