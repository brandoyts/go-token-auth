package user

type Handler struct {
	userService *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{userService: userService}
}
