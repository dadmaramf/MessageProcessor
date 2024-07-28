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
	const op = "storage.postgres.PostMessage"

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

	query, err := tx.Prepare("INSERT INTO message (content, status) VALUES($1, 'none')")
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

	query, err := tx.Prepare("INSERT INTO outbox (content, status) VALUES($1, $2)")
	if err != nil {
		return err
	}
	_, err = query.Exec(msg, "new")

	return err
}

type MessageOutbox struct {
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
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	row := query.QueryRowContext(ctx, time.Now())

	var outbox MessageOutbox

	err = row.Scan(&outbox.ID, &outbox.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil, nil
		}
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if err = m.updateOutbox(ctx, tx, outbox.ID); err != nil {
		return nil, err
	}

	if errCommit := tx.Commit(); errCommit != nil {
		err = fmt.Errorf("%s:%w", op, errCommit)
		return
	}

	return &model.Message{
		ID:      outbox.ID,
		Content: outbox.Content,
	}, nil

}

func (m *MessageStorage) updateOutbox(ctx context.Context, tx *sql.Tx, id int) error {
	const op = "storage.postgres.updateOutbox"
	query, err := tx.PrepareContext(ctx, "UPDATE outbox SET reserved = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = query.Exec(time.Now().Add(reserved), id)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return err
}

// SetDown updates the status of a message in the outbox to 'done'.
func (m *MessageStorage) SetDown(id int) error {
	const op = "storage.postgres.SetDown"

	query, err := m.db.Prepare("UPDATE outbox SET status = 'done' WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = query.Exec(id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	return nil
}

type Message struct {
	ID               int    `json:"id" db:"id"`
	Content          string `json:"content" db:"content"`
	ProcessedContent string `json:"processed_content" db:"processed_content"`
	Status           string `json:"status" db:"status"`
}

// GetDownMessages retrieves messages from the outbox with status 'done'.
func (m *MessageStorage) GetDownMessages() ([]model.Message, error) {
	const op = "storage.postgres.SetDown"

	query, err := m.db.Query("SELECT id, content, processed_content, status FROM message WHERE status = 'update'")
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer query.Close()

	message := make([]model.Message, 0, 100)
	for query.Next() {
		var msg Message
		err = query.Scan(&msg.ID, &msg.Content, &msg.ProcessedContent, &msg.Status)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}
		message = append(message, model.Message{ID: msg.ID, Content: msg.Content, ProcessedContent: msg.ProcessedContent, Status: msg.Status})
	}

	if query.Err() != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return message, nil
}

func (m *MessageStorage) AddProcessedMessage(id int, msg string) error {
	const op = "storage.postgres.AddProcessedMessage"

	query, err := m.db.Prepare("UPDATE message SET status = 'update', processed_content = $1 WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = query.Exec(msg, id)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
