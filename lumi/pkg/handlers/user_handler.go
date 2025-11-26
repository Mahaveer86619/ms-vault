package handlers

import (
	"net/http"

	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/views"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(group *echo.Group, us *services.UserService) *UserHandler {
	handler := &UserHandler{
		userService: us,
	}

	group.GET("/me", handler.GetUserDetails)

	return handler
}

func (h *UserHandler) GetUserDetails(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusUnauthorized)
		failure.SetMessage("Unauthorized")
		return failure.JSON(c)
	}

	resp, err := h.userService.GetUserDetails(userID)
	if err != nil {
		failure := views.Failure{}
		failure.SetStatusCode(http.StatusInternalServerError)
		failure.SetMessage(err.Error())
		return failure.JSON(c)
	}

	success := views.Success{}
	success.SetStatusCode(http.StatusOK)
	success.SetMessage("User details fetched successfully")
	success.SetData(resp)

	return success.JSON(c)
}
