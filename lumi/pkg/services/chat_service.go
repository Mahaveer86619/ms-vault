package services

import (
	"log"

	"github.com/Mahaveer86619/lumi/pkg/db"
	"github.com/Mahaveer86619/lumi/pkg/models"
	connModel "github.com/Mahaveer86619/lumi/pkg/models/connections"
	"github.com/Mahaveer86619/lumi/pkg/services/connections"
)

type ChatService struct {
	WahaClient connections.WahaClient
}

func NewChatService(wahaClient connections.WahaClient) *ChatService {
	return &ChatService{
		WahaClient: wahaClient,
	}
}

func (s *ChatService) GetRemoteChats() ([]connModel.ChatSummary, error) {
	return s.WahaClient.GetChats()
}

func (s *ChatService) GetRemoteGroups() ([]connModel.GroupInfo, error) {
	return s.WahaClient.GetGroups()
}

func (s *ChatService) GetRegisteredChats() ([]models.RegisteredChat, error) {
	var chats []models.RegisteredChat
	result := db.DB.Find(&chats)

	if len(chats) == 0 {
		log.Println("0 chats registered")
		chats = []models.RegisteredChat{}
	}
	return chats, result.Error
}

func (s *ChatService) RegisterChat(chatID, name, chatType string) (*models.RegisteredChat, error) {
	chat := models.RegisteredChat{
		ChatID: chatID,
		Name:   name,
		Type:   chatType,
	}

	if err := db.DB.Where("chat_id = ?", chatID).FirstOrCreate(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (s *ChatService) UnregisterChat(chatID string) error {
	return db.DB.Where("chat_id = ?", chatID).Unscoped().Delete(&models.RegisteredChat{}).Error
}

func (s *ChatService) IsChatAllowed(chatID string) bool {
	var count int64
	db.DB.Model(&models.RegisteredChat{}).Where("chat_id = ?", chatID).Count(&count)
	return count > 0
}
