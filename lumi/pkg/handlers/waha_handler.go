package handlers

import (
	"fmt"
	"net/http"

	"github.com/Mahaveer86619/lumi/pkg/services"
	"github.com/Mahaveer86619/lumi/pkg/views"
	"github.com/labstack/echo/v4"
)

type WahaHandler struct {
	wahaService *services.WahaService
}

func NewWahaHandler(group *echo.Group, wahaService *services.WahaService) *WahaHandler {
	handler := &WahaHandler{
		wahaService: wahaService,
	}

	group.GET("/connect", handler.ConnectWhatsApp)
	group.GET("/me", handler.ConnectWhatsApp)

	return handler
}

func (h *WahaHandler) ConnectWhatsApp(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, views.Failure{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized: User ID not found",
		})
	}

	sessionName := fmt.Sprintf("user_%d", userID)

	err := h.wahaService.StartSession(sessionName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, views.Failure{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to start WhatsApp session: " + err.Error(),
		})
	}

	profile, err := h.wahaService.GetProfile(sessionName)
	if err == nil && profile != nil {
		return c.JSON(http.StatusOK, views.Success{
			StatusCode: http.StatusOK,
			Message:    "Already logged in",
			Data: map[string]interface{}{
				"status":  "connected",
				"profile": profile,
			},
		})
	}

	qrBytes, err := h.wahaService.GetQRCode(sessionName)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.Failure{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to retrieve QR code from Waha: " + err.Error(),
		})
	}

	return c.Blob(http.StatusOK, "image/png", qrBytes)
}

func (h *WahaHandler) GetMe(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, views.Failure{
			StatusCode: http.StatusUnauthorized,
			Message:    "Unauthorized",
		})
	}

	sessionName := fmt.Sprintf("user_%d", userID)

	profile, err := h.wahaService.GetProfile(sessionName)
	if err != nil {
		return c.JSON(http.StatusBadGateway, views.Failure{
			StatusCode: http.StatusBadGateway,
			Message:    "Failed to fetch profile. Ensure session is connected.",
		})
	}

	return c.JSON(http.StatusOK, views.Success{
		StatusCode: http.StatusOK,
		Message:    "Profile fetched successfully",
		Data:       profile,
	})
}
