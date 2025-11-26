package views

type HealthResponse struct {
	Services []Health `json:"services"`
}

type Health struct {
	Name    string `json:"name"`
	IsUp    bool   `json:"is_up"`
	Message string `json:"message"`
}
