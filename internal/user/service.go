package user

type Service struct {
	repository DBInterface
}

func NewService(repo DBInterface) *Service {
	return &Service{repository: repo}
}
