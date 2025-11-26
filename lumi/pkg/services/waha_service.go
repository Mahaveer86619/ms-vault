package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
	"github.com/Mahaveer86619/lumi/pkg/enums"
	"github.com/Mahaveer86619/lumi/pkg/models"
)

type WahaService struct {
	httpClient *http.Client
}

func NewWahaService() *WahaService {
	return &WahaService{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (s *WahaService) StartSession(sessionName string) error {
	status, err := s.getSessionStatus(sessionName)
	if err != nil {
		if err.Error() == "session not found" {
			if err := s.createSession(sessionName); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if status == enums.WAHA_SESSION_STOPPED.String() || status == enums.WAHA_SESSION_FAILED.String() {
			if err := s.startExistingSession(sessionName); err != nil {
				return err
			}
		}
	}

	return s.waitForSessionReady(sessionName)
}

func (s *WahaService) GetQRCode(sessionName string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s/auth/qr?format=image", config.GConfig.WahaServiceURL, sessionName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get QR code, status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

func (s *WahaService) GetProfile(sessionName string) (*models.WahaProfile, error) {
	url := fmt.Sprintf("%s/api/%s/profile", config.GConfig.WahaServiceURL, sessionName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get profile: %d %s", resp.StatusCode, string(body))
	}

	var profile models.WahaProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (s *WahaService) waitForSessionReady(sessionName string) error {
	timeout := time.After(20 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for session %s to be ready", sessionName)
		case <-ticker.C:
			status, err := s.getSessionStatus(sessionName)
			if err != nil {
				continue
			}

			if status == enums.WAHA_SESSION_SCAN_QR_CODE.String() || status == enums.WAHA_SESSION_WORKING.String() {
				return nil
			}

			if status == enums.WAHA_SESSION_FAILED.String() {
				return fmt.Errorf("session failed to start")
			}

			if status == enums.WAHA_SESSION_STOPPED.String() {
				return fmt.Errorf("session stopped unexpectedly")
			}
		}
	}
}

func (s *WahaService) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", config.GConfig.WahaAPIKey)
}

func (s *WahaService) getSessionStatus(sessionName string) (string, error) {
	url := fmt.Sprintf("%s/api/sessions/%s", config.GConfig.WahaServiceURL, sessionName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("session not found")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get session: %d %s", resp.StatusCode, string(body))
	}

	var info struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", err
	}
	return info.Status, nil
}

func (s *WahaService) createSession(sessionName string) error {
	url := fmt.Sprintf("%s/api/sessions", config.GConfig.WahaServiceURL)
	payload := map[string]interface{}{
		"name":  sessionName,
		"start": true,
	}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create session: %d %s", resp.StatusCode, string(body))
	}
	return nil
}

func (s *WahaService) startExistingSession(sessionName string) error {
	url := fmt.Sprintf("%s/api/sessions/%s/start", config.GConfig.WahaServiceURL, sessionName)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to start session: %d %s", resp.StatusCode, string(body))
	}
	return nil
}
