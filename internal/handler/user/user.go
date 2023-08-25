package user

import (
	"avito-internship-2023/internal/api"
	"avito-internship-2023/internal/service"
)

type Handler struct {
	userService service.IUser
	r           *api.Renderer
}

func NewHandler(
	userService service.IUser,
	renderer *api.Renderer,
) *Handler {
	return &Handler{
		userService: userService,
		r:           renderer,
	}
}
