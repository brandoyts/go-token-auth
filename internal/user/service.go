package user

import (
	"context"
)

type Service struct {
	repository DBInterface
}

func NewService(repo DBInterface) *Service {
	return &Service{repository: repo}
}

func (s *Service) FindUserById(ctx context.Context, id string) (*User, error) {
	result, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) FindUser(ctx context.Context, user User) (*User, error) {
	result, err := s.repository.FindOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Service) CreateUser(ctx context.Context, user User) (string, error) {
	err := user.HashPassword()
	if err != nil {
		return "", err
	}

	result, err := s.repository.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return result, nil
}
