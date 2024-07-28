package model

type Message struct {
	ID               int    `json:"id" db:"id"`
	Content          string `json:"content" db:"content"`
	ProcessedContent string `json:"processed_content,omitempty" db:"processed_content"`
	Status           string `json:"status,omitempty" db:"status"`
}

type MessageState struct {
	ID       int    `json:"id" db:"id"`
	Content  string `json:"content" db:"content"`
	Status   string `json:"status" db:"status"`
	CreateAt string `json:"cerate_at" db:"cerate_at"`
}
