package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/ms/auth/pkg/services"
	"github.com/Mahaveer86619/ms/auth/pkg/views"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(group *echo.Group, authService *services.AuthService) *AuthHandler {
	handler := &AuthHandler{
		authService: authService,
	}

	// Register the routes
	group.POST("/register", handler.RegisterWithEmail)
	group.POST("/login", handler.LoginWithEmail)

	return handler
}

func (h *AuthHandler) RegisterWithEmail(c echo.Context) error {
	
	var req views.AuthRegisterWithEmail
	if err := c.Bind(&req); err != nil {
		failure := &views.Failure{}
		failure.SetStatusCode(http.StatusBadRequest)
		failure.SetMessage("Invalid request payload")
		return failure.JSON(c)
	}

	resp, err := h.authService.RegisterWithEmail(req.Username, req.Email, req.Password)
	if err != nil {
		failure := &views.Failure{}
		failure.SetStatusCode(http.StatusInternalServerError)
		failure.SetMessage("Registration failed")
		return failure.JSON(c)
	}

	success := &views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("Registration successful")
	success.SetData(resp)

	return success.JSON(c)
	
}

func (h *AuthHandler) LoginWithEmail(c echo.Context) error {
	return c.String(http.StatusOK, "Login logic goes here. Needs JSON response.")
}
