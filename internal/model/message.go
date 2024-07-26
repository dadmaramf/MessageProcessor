package model

type Message struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Status  string `json:"status" db:"status"`
}
