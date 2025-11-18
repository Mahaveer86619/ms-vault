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

	group.POST("/register", handler.RegisterUser)
	group.POST("/login", handler.LoginWithUsername)
	group.POST("/refresh", handler.RefreshToken)

	return handler
}

func (h *AuthHandler) RegisterUser(c echo.Context) error {

	var req views.AuthRegisterUser
	if err := c.Bind(&req); err != nil {
		failure := &views.Failure{}
		failure.SetStatusCode(http.StatusBadRequest)
		failure.SetMessage("Invalid request payload")
		return failure.JSON(c)
	}

	resp, err := h.authService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		if err.Error() == "username already exists" {
			failure := &views.Failure{}
			failure.SetStatusCode(http.StatusBadRequest)
			failure.SetMessage(err.Error())
			return failure.JSON(c)
		}
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

func (h *AuthHandler) LoginWithUsername(c echo.Context) error {

	var req views.AuthLoginWithUsername
	if err := c.Bind(&req); err != nil {
		failure := &views.Failure{}
		failure.SetStatusCode(http.StatusBadRequest)
		failure.SetMessage("Invalid request payload")
		return failure.JSON(c)
	}

	resp, err := h.authService.LoginWithUsername(req.Username, req.Password)
	if err != nil {
		failure := &views.Failure{}
		failure.SetStatusCode(http.StatusUnauthorized)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	success := &views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("login successful")
	success.SetData(resp)

	return success.JSON(c)
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	
	var req views.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		failure := &views.Failure{}
		failure.SetStatusCode(http.StatusBadRequest)
		failure.SetMessage("Invalid request payload")
		return failure.JSON(c)
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		failure := &views.Failure{}
		failure.SetStatusCode(http.StatusUnauthorized)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	success := &views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("Token refreshed successfully")
	success.SetData(resp)

	return success.JSON(c)
}
