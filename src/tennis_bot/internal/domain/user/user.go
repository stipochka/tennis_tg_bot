package user

import "time"

type User struct {
	ID         int64
	TelegramID int64 `json:"telegram_id"`
	IsAdmin    bool
	CreatedAt  time.Time
}
