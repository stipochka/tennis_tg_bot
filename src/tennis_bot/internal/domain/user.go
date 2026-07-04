package domain

type User struct {
	ID         int64
	TelegramID string `json:"telegram_id"`
	IsAdmin    bool
}
