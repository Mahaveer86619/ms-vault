package bot

import (
	"log"
	"strings"

	modelConnections "github.com/Mahaveer86619/lumi/pkg/models/connections"
	"github.com/Mahaveer86619/lumi/pkg/services/connections"
)

type BotService struct {
	wahaClient connections.WahaClient
}

func NewBotService(wahaClient connections.WahaClient) *BotService {
	return &BotService{
		wahaClient: wahaClient,
	}
}

func (b *BotService) ProcessMessage(msg modelConnections.WAMessage) {
	// 1. Sanity Checks
	if msg.FromMe {
		return
	}

	chatID := msg.From
	text := strings.TrimSpace(msg.Body)

	// 2. Command Parsing (Basic Implementation)
	if strings.HasPrefix(text, "/") {
		b.handleCommand(chatID, text)
		return
	}

	// 3. Default Behavior (e.g., Echo or NLP processing)
	// For now, we keep the echo logic here
	b.sendEcho(chatID, text)
}

func (b *BotService) handleCommand(chatID, text string) {
	parts := strings.SplitN(text, " ", 2)
	command := parts[0]

	switch command {
	case "/ping":
		b.wahaClient.SendText(chatID, "Pong! 🏓")
	case "/help":
		b.wahaClient.SendText(chatID, "Available commands: /ping, /help")
	default:
		b.wahaClient.SendText(chatID, "Unknown command.")
	}
}

func (b *BotService) sendEcho(chatID, text string) {
	if text == "" {
		return
	}
	_, err := b.wahaClient.SendText(chatID, "You said: "+text)
	if err != nil {
		log.Printf("Bot failed to reply: %v", err)
	}
}
