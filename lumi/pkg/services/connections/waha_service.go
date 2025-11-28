package connections

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Mahaveer86619/lumi/pkg/config"
	models "github.com/Mahaveer86619/lumi/pkg/models/connections"
)

type WahaClient interface {
	// Lifecycle
	Ping() error
	StartSession() error
	StopSession() error
	RestartSession() error
	GetSessionStatus() (*models.SessionInfo, error)
	GetQRCode() ([]byte, error)
	RequestCode(phoneNumber string, method string) (*models.RequestCodeResponse, error)
	GetMe() (*models.MeInfo, error)

	// Chatting
	SendText(chatId, text string) (*models.WAMessage, error)
	SendImage(chatId string, image models.ImagePayload) (*models.WAMessage, error)
	CheckNumberExists(phone string) (*models.WANumberExistResult, error)
	GetChats() ([]models.ChatSummary, error)
	GetGroups() ([]models.GroupInfo, error)
}

type WahaService struct {
	httpClient  *http.Client
	baseURL     string
	sessionName string
	apiKey      string
}

func NewWahaService() WahaClient {
	return &WahaService{
		httpClient:  &http.Client{Timeout: 60 * time.Second},
		baseURL:     config.GConfig.WahaServiceURL,
		sessionName: config.GConfig.WahaSessionName,
		apiKey:      config.GConfig.WahaAPIKey,
	}
}

// --- Lifecycle Methods ---

func (s *WahaService) Ping() error {
	url := fmt.Sprintf("%s/ping", s.baseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	return s.doRequest(req, nil)
}

func (s *WahaService) GetSessionStatus() (*models.SessionInfo, error) {
	url := fmt.Sprintf("%s/api/sessions/%s", s.baseURL, s.sessionName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var sessionInfo models.SessionInfo
	if err := s.doRequest(req, &sessionInfo); err != nil {
		return nil, err
	}
	return &sessionInfo, nil
}

func (s *WahaService) StartSession() error {
	info, err := s.GetSessionStatus()

	if err != nil {
		if err := s.createSession(); err != nil {
			return fmt.Errorf("failed to create session after status check failed: %w", err)
		}
	} else if info != nil {
		switch info.Status {
		case "STOPPED", "FAILED":
			if err := s.startExistingSession(); err != nil {
				return err
			}
		case "WORKING", "SCAN_QR_CODE", "STARTING":
			// Already running or starting, just wait
		}
	}

	return s.waitForSessionReady()
}

func (s *WahaService) createSession() error {
	url := fmt.Sprintf("%s/api/sessions", s.baseURL)

	payload := models.SessionCreateRequest{
		Name:  s.sessionName,
		Start: true,
		// Config: &models.SessionConfig{ ... }
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	return s.doRequest(req, nil)
}

func (s *WahaService) startExistingSession() error {
	url := fmt.Sprintf("%s/api/sessions/%s/start", s.baseURL, s.sessionName)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	return s.doRequest(req, nil)
}

func (s *WahaService) StopSession() error {
	url := fmt.Sprintf("%s/api/sessions/%s/stop", s.baseURL, s.sessionName)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	return s.doRequest(req, nil)
}

func (s *WahaService) RestartSession() error {
	url := fmt.Sprintf("%s/api/sessions/%s/restart", s.baseURL, s.sessionName)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	s.doRequest(req, nil)

	return s.waitForSessionReady()
}

func (s *WahaService) GetQRCode() ([]byte, error) {
	url := fmt.Sprintf("%s/api/%s/auth/qr?format=image", s.baseURL, s.sessionName)
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
		return nil, fmt.Errorf("failed to get QR code: %d %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (s *WahaService) RequestCode(phoneNumber string, method string) (*models.RequestCodeResponse, error) {
	url := fmt.Sprintf("%s/api/%s/auth/request-code", s.baseURL, s.sessionName)

	payload := make(map[string]string)
	payload["phoneNumber"] = phoneNumber

	if method != "" {
		payload["method"] = method
	}

	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		jsonPayload, _ := json.Marshal(payload)
		req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if reqErr != nil {
			return nil, reqErr
		}

		var response models.RequestCodeResponse
		err = s.doRequest(req, &response)

		if err == nil {
			return &response, nil
		}

		if i < maxRetries-1 {
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
}

func (s *WahaService) GetMe() (*models.MeInfo, error) {
	url := fmt.Sprintf("%s/api/sessions/%s/me", s.baseURL, s.sessionName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var meInfo models.MeInfo
	if err := s.doRequest(req, &meInfo); err != nil {
		return nil, err
	}
	return &meInfo, nil
}

// --- Chatting Methods ---

func (s *WahaService) SendText(chatId, text string) (*models.WAMessage, error) {
	url := fmt.Sprintf("%s/api/sendText", s.baseURL)

	payload := models.MessageTextRequest{
		ChatID:  chatId,
		Text:    text,
		Session: s.sessionName,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	var response models.WAMessage
	if err := s.doRequest(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *WahaService) SendImage(chatId string, image models.ImagePayload) (*models.WAMessage, error) {
	url := fmt.Sprintf("%s/api/sendImage", s.baseURL)

	payload := models.MessageImageRequest{
		ChatID:  chatId,
		Session: s.sessionName,
		Caption: image.Caption,
		File:    image.File,
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	var response models.WAMessage
	if err := s.doRequest(req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (s *WahaService) CheckNumberExists(phone string) (*models.WANumberExistResult, error) {
	url := fmt.Sprintf("%s/api/contacts/check-exists?phone=%s&session=%s", s.baseURL, phone, s.sessionName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result models.WANumberExistResult
	if err := s.doRequest(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *WahaService) GetChats() ([]models.ChatSummary, error) {
	url := fmt.Sprintf("%s/api/%s/chats/overview", s.baseURL, s.sessionName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var chats []models.ChatSummary
	if err := s.doRequest(req, &chats); err != nil {
		return nil, err
	}
	return chats, nil
}

func (s *WahaService) GetGroups() ([]models.GroupInfo, error) {
	url := fmt.Sprintf("%s/api/%s/groups", s.baseURL, s.sessionName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var groups []models.GroupInfo
	if err := s.doRequest(req, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

// --- Helpers ---

func (s *WahaService) addHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", s.apiKey)
}

func (s *WahaService) doRequest(req *http.Request, v interface{}) error {
	s.addHeaders(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

func (s *WahaService) waitForSessionReady() error {
	timeout := time.After(20 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for session %s to be ready", s.sessionName)
		case <-ticker.C:
			status, err := s.GetSessionStatus()
			if err != nil {
				continue
			}

			if status.Status == "SCAN_QR_CODE" || status.Status == "WORKING" {
				return nil
			}
			if status.Status == "FAILED" {
				return fmt.Errorf("session failed to start")
			}
			if status.Status == "STOPPED" {
				continue
			}
		}
	}
}
