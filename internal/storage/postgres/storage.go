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

var reserved = time.Hour

type MessageStorage struct {
	db *sql.DB
}

func NewMessageStorage(db *sql.DB) *MessageStorage {
	return &MessageStorage{
		db: db,
	}
}

func (m *MessageStorage) PostMessage(msg string) (err error) {
	const op = "storage.postgres.Post"
	tx, err := m.db.Begin()

	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}

		errCommit := tx.Commit()
		if errCommit != nil {
			err = fmt.Errorf("%s:%v", op, errCommit)
		}
	}()

	query, err := tx.Prepare("INSERT INTO message (content) VALUES(?)")

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

func (m *MessageStorage) postOutbox(tx *sql.Tx, msg string) error {

	query, err := tx.Prepare("INSERT INTO outbox (content, status, cerate_at, reserved) VALUES(?, ?, ?, ?)")

	if err != nil {
		return err
	}
	_, err = query.Exec(msg, "new", time.Now(), time.Now().Add(reserved))

	return err
}

type Message struct {
	ID      int    `json:"id" db:"id"`
	Content string `json:"content" db:"content"`
	Status  string `json:"status" db:"status"`
}

func (m *MessageStorage) GetNewOutbox(ctx context.Context) (*model.Message, error) {
	const op = "storage.postgres.GetNewOutbox"

	query, err := m.db.PrepareContext(ctx, "SELECT id, content FROM outbox WHERE status = 'new' AND reserved < ? LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}
	row := query.QueryRowContext(ctx, time.Now())

	var outbox Message

	err = row.Scan(&outbox.ID, &outbox.Content)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	return &model.Message{
		ID:      outbox.ID,
		Content: outbox.Content,
	}, nil

}

func (m *MessageStorage) SetDown(id int) error {
	const op = "storage.postgres.SetDown"

	query, err := m.db.Prepare("UPDATE outbox SET status = 'down' WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}

	_, err = query.Exec(id)

	if err != nil {
		return fmt.Errorf("%s:%v", op, err)
	}
	return nil
}

func (m *MessageStorage) GetDownMessages() ([]model.Message, error) {
	const op = "storage.postgres.SetDown"

	query, err := m.db.Query("SELECT id, content FROM outbox WHERE status = 'down'")

	if err != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	defer query.Close()
	message := make([]model.Message, 0, 100)

	for query.Next() {
		var msg Message
		err = query.Scan(&msg.ID, &msg.Content)
		if err != nil {
			return nil, fmt.Errorf("%s:%v", op, err)
		}

		message = append(message, model.Message{ID: msg.ID, Content: msg.Content})
	}

	if query.Err() != nil {
		return nil, fmt.Errorf("%s:%v", op, err)
	}

	return message, nil
}
