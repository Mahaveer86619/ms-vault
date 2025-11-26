package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/views"
)

type HealthService struct {
	httpClient *http.Client
}

func NewHealthService() *HealthService {
	return &HealthService{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *HealthService) GetHealth() (*views.HealthResponse, error) {
	var servicesList []views.Health

	// 1. Check Waha Service
	servicesList = append(servicesList, h.checkWahaService())

	// 2. Check Database
	servicesList = append(servicesList, h.checkDBService())

	// 3. Check Lumi Service
	servicesList = append(servicesList, views.Health{
		Name:    "lumi-service",
		IsUp:    true,
		Message: "Service is running",
	})

	return &views.HealthResponse{
		Services: servicesList,
	}, nil
}

func (h *HealthService) checkWahaService() views.Health {
	url := fmt.Sprintf("%s/api/sessions", config.GConfig.WahaServiceURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return views.Health{
			Name:    "waha-service",
			IsUp:    false,
			Message: fmt.Sprintf("Request creation failed: %v", err),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", config.GConfig.WahaAPIKey)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return views.Health{
			Name:    "waha-service",
			IsUp:    false,
			Message: fmt.Sprintf("Connection failed: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return views.Health{
			Name:    "waha-service",
			IsUp:    false,
			Message: fmt.Sprintf("Unhealthy status code: %d", resp.StatusCode),
		}
	}

	return views.Health{
		Name:    "waha-service",
		IsUp:    true,
		Message: "Healthy",
	}
}

func (h *HealthService) checkDBService() views.Health {
	if db.DB == nil {
		return views.Health{
			Name:    "database",
			IsUp:    false,
			Message: "Database connection not initialized",
		}
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return views.Health{
			Name:    "database",
			IsUp:    false,
			Message: fmt.Sprintf("Failed to retrieve DB instance: %v", err),
		}
	}

	if err := sqlDB.Ping(); err != nil {
		return views.Health{
			Name:    "database",
			IsUp:    false,
			Message: fmt.Sprintf("Ping failed: %v", err),
		}
	}

	return views.Health{
		Name:    "database",
		IsUp:    true,
		Message: "Healthy",
	}
}
