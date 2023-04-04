package entity

import "time"

type Message struct {
	ID        int       `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	UserID    int       `json:"user_id" db:"user_id"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
