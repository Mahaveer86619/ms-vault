package views

type RemoteChatListResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	LastMessage string `json:"last_message"`
	Timestamp   int64  `json:"timestamp"`
	Type        string `json:"type"`
}

type RegisterChatRequest struct {
	ChatID string `json:"chat_id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

type SendTextChatRequest struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}
