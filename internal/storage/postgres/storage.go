package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messageprocessor/internal/model"
	"messageprocessor/internal/storage"
	"time"
)

var _ storage.Storage = (*MessageStorage)(nil)

const reserved = time.Hour

type MessageStorage struct {
	db *sql.DB
}

func NewMessageStorage(db *sql.DB) *MessageStorage {
	return &MessageStorage{
		db: db,
	}
}

// PostMessage saves a message to the database.
func (m *MessageStorage) PostMessage(msg string) (err error) {
	const op = "storage.postgres.Post"

	if msg == "" {
		return fmt.Errorf("%s: message content is empty", op)
	}

	tx, err := m.db.Begin()

	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		errCommit := tx.Commit()
		if errCommit != nil {
			err = fmt.Errorf("%s:%v", op, errCommit)
		}
	}()

	query, err := tx.Prepare("INSERT INTO message (content) VALUES($1)")

	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = query.Exec(msg)

	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	err = m.postOutbox(tx, msg)

	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	return nil
}

// postOutbox adds a message to the outbox table.
func (m *MessageStorage) postOutbox(tx *sql.Tx, msg string) error {

	query, err := tx.Prepare("INSERT INTO outbox (content, status, create_at) VALUES($1, $2, $3)")

	if err != nil {
		return err
	}
	_, err = query.Exec(msg, "new", time.Now())

	return err
}

type Message struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Status  string `json:"status" db:"status"`
}

// GetNewOutbox retrieves a new message from the outbox table that is ready to be sent.
func (m *MessageStorage) GetNewOutbox(ctx context.Context) (msg *model.Message, err error) {
	const op = "storage.postgres.GetNewOutbox"

	tx, err := m.db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query, err := tx.PrepareContext(ctx, "SELECT id, content FROM outbox WHERE status = 'new' AND (reserved IS NULL OR reserved < $1) LIMIT 1")

	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	row := query.QueryRowContext(ctx, time.Now())

	var outbox Message

	err = row.Scan(&outbox.ID, &outbox.Content)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil, nil
		}
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	query, err = tx.PrepareContext(ctx, "UPDATE outbox SET reserved = $1 WHERE id = $2")
	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}
	_, err = query.Exec(time.Now().Add(reserved), outbox.ID)

	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	if errCommit := tx.Commit(); errCommit != nil {
		err = fmt.Errorf("%s:%v", op, errCommit)
		return
	}

	return &model.Message{
		ID:      outbox.ID,
		Content: outbox.Content,
	}, nil

}

// SetDown updates the status of a message in the outbox to 'done'.
func (m *MessageStorage) SetDown(id int) error {
	const op = "storage.postgres.SetDown"

	query, err := m.db.Prepare("UPDATE outbox SET status = 'done' WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = query.Exec(id)

	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}
	return nil
}

type MessageState struct {
	ID       int    `json:"id" db:"id"`
	Content  string `json:"content" db:"content"`
	Status   string `json:"status" db:"status"`
	CreateAt string `json:"cerate_at" db:"cerate_at"`
}

// GetDownMessages retrieves messages from the outbox with status 'done'.
func (m *MessageStorage) GetDownMessages() ([]model.MessageState, error) {
	const op = "storage.postgres.SetDown"

	query, err := m.db.Query("SELECT id, content, create_at, status FROM outbox WHERE status = 'done'")

	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	defer query.Close()
	message := make([]model.MessageState, 0, 100)

	for query.Next() {
		var msg MessageState
		err = query.Scan(&msg.ID, &msg.Content, &msg.CreateAt, &msg.Status)
		if err != nil {
			return nil, fmt.Errorf("%s:%v", op, err)
		}

		message = append(message, model.MessageState{ID: msg.ID, Content: msg.Content, CreateAt: msg.CreateAt, Status: msg.Status})
	}

	if query.Err() != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	return message, nil
}
