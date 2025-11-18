package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/ms/auth/pkg/services"
	"github.com/labstack/echo/v4"
)

type AvatarHandler struct {
	avatarService *services.AvatarService
}

func NewAvatarHandler(group *echo.Group, avatarService *services.AvatarService) *AvatarHandler {
	handler := &AvatarHandler{
		avatarService: avatarService,
	}

	group.GET("/avatar/:hash", handler.GetAvatar)

	return handler
}

func (h *AvatarHandler) GetAvatar(c echo.Context) error {
	hash := c.Param("hash")
	if hash == "" {
		return c.String(http.StatusNotFound, "Avatar hash not provided")
	}

	img, err := h.avatarService.GenerateAvatarImage(hash)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to generate avatar image")
	}

	imgBytes, err := h.avatarService.ImageToBytes(img)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to encode image")
	}

	c.Response().Header().Set("Content-Type", "image/png")
	c.Response().WriteHeader(http.StatusOK)
	_, writeErr := c.Response().Write(imgBytes)

	return writeErr
}
