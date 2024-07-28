package storage_test

import (
	"context"
	"database/sql"
	"messageprocessor/internal/model"
	post "messageprocessor/internal/storage/postgres"
	storage "messageprocessor/internal/storage/postgres"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	service := post.NewMessageStorage(db)

	tests := []struct {
		argfirstQuery     string
		returnfirstQuery  []int64
		argsecondQuery    []any
		returnsecondQuery []int64
	}{
		{
			argfirstQuery:     "love",
			returnfirstQuery:  []int64{1, 1},
			argsecondQuery:    []any{"love", "new"},
			returnsecondQuery: []int64{1, 1},
		},
		{
			argfirstQuery:     "hi",
			returnfirstQuery:  []int64{1, 1},
			argsecondQuery:    []any{"hi", "new"},
			returnsecondQuery: []int64{1, 1},
		},
		{
			argfirstQuery:     "",
			returnfirstQuery:  []int64{1, 1},
			argsecondQuery:    []any{"", "new"},
			returnsecondQuery: []int64{1, 1},
		},
	}
	for _, test := range tests {

		if test.argfirstQuery == "" {
			if err = service.PostMessage(test.argfirstQuery); err == nil {
				t.Errorf("expected error for empty message, but got none")
			}
			continue
		}

		mock.ExpectBegin()
		mock.ExpectPrepare("INSERT INTO message \\(content, status\\) VALUES\\(\\$1, 'none'\\)").ExpectExec().WithArgs(test.argfirstQuery).WillReturnResult(sqlmock.NewResult(test.returnfirstQuery[0], test.returnfirstQuery[1]))
		mock.ExpectPrepare("INSERT INTO outbox \\(content, status\\) VALUES\\(\\$1, \\$2\\)").ExpectExec().WithArgs(test.argsecondQuery[0], test.argsecondQuery[1]).WillReturnResult(sqlmock.NewResult(test.returnsecondQuery[0], test.returnsecondQuery[1]))
		mock.ExpectCommit()

		err = service.PostMessage(test.argfirstQuery)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}

}

func TestGetNewOutbox(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := post.NewMessageStorage(db)

	tests := []struct {
		Name            string
		ExpectedMessage *model.Message
		ArgSecondQuery  int
		PrepareError    error
		QueryError      error
		ExecError       error
		ExpectedError   error
	}{
		{
			Name:            "Success",
			ExpectedMessage: &model.Message{ID: 1, Content: "test message"},
			ArgSecondQuery:  1,
			PrepareError:    nil,
			QueryError:      nil,
			ExecError:       nil,
			ExpectedError:   nil,
		},
		{
			Name:            "NoRows",
			ExpectedMessage: nil,
			ArgSecondQuery:  1,
			PrepareError:    nil,
			QueryError:      sql.ErrNoRows,
			ExecError:       nil,
			ExpectedError:   nil,
		},
		{
			Name:            "ErrorPrepareSelect",
			ExpectedMessage: nil,
			ArgSecondQuery:  1,
			PrepareError:    sql.ErrConnDone,
			QueryError:      nil,
			ExecError:       nil,
			ExpectedError:   sql.ErrConnDone,
		},
		{
			Name:            "ErrorPrepareUpdate",
			ExpectedMessage: &model.Message{ID: 1, Content: "test message"},
			ArgSecondQuery:  1,
			PrepareError:    nil,
			QueryError:      nil,
			ExecError:       sql.ErrConnDone,
			ExpectedError:   sql.ErrConnDone,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			mock.ExpectBegin()

			if test.PrepareError != nil {
				mock.ExpectPrepare("SELECT id, content FROM outbox WHERE status = 'new' AND \\(reserved IS NULL OR reserved < \\$1\\) LIMIT 1").
					WillReturnError(test.PrepareError)
				mock.ExpectRollback()
			} else {
				rows := sqlmock.NewRows([]string{"id", "content"})
				if test.ExpectedMessage != nil {
					rows.AddRow(test.ExpectedMessage.ID, test.ExpectedMessage.Content)
				}
				mock.ExpectPrepare("SELECT id, content FROM outbox WHERE status = 'new' AND \\(reserved IS NULL OR reserved < \\$1\\) LIMIT 1").
					ExpectQuery().
					WithArgs(sqlmock.AnyArg()). // use sqlmock.AnyArg() to match any argument
					WillReturnRows(rows).
					WillReturnError(test.QueryError)

				if test.QueryError == nil && test.ExpectedMessage != nil {
					if test.ExecError != nil {
						mock.ExpectPrepare("UPDATE outbox SET reserved = \\$1 WHERE id = \\$2").
							ExpectExec().
							WithArgs(sqlmock.AnyArg(), test.ArgSecondQuery).
							WillReturnError(test.ExecError)
						mock.ExpectRollback()
					} else {
						mock.ExpectPrepare("UPDATE outbox SET reserved = \\$1 WHERE id = \\$2").
							ExpectExec().
							WithArgs(sqlmock.AnyArg(), test.ArgSecondQuery).
							WillReturnResult(sqlmock.NewResult(1, 1))
						mock.ExpectCommit()
					}
				} else if test.QueryError == sql.ErrNoRows {
					mock.ExpectRollback()
				}
			}

			msgRes, err := service.GetNewOutbox(context.Background())

			if err != nil && (test.ExpectedError == nil) {
				t.Errorf("%s: expected error '%v', got '%v'", test.Name, test.ExpectedError, err)
			}

			if err == nil && test.ExpectedError != nil {
				t.Errorf("%s: expected error '%v', but got nil", test.Name, test.ExpectedError)
			}

			if err == nil && msgRes == nil && test.ExpectedMessage != nil {
				t.Errorf("%s: expected message '%v', got nil", test.Name, test.ExpectedMessage)
			}

			if err == nil && test.ExpectedMessage != nil && (msgRes.ID != test.ExpectedMessage.ID || msgRes.Content != test.ExpectedMessage.Content) {
				t.Errorf("%s: expected message '%v', got '%v'", test.Name, test.ExpectedMessage, msgRes)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestSetDown(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	service := post.NewMessageStorage(db)

	tests := []struct {
		Name          string
		ID            int
		PrepareError  error
		ExecError     error
		ExpectedError error
	}{
		{
			Name:          "Success",
			ID:            1,
			PrepareError:  nil,
			ExecError:     nil,
			ExpectedError: nil,
		},
		{
			Name:          "ErrorPrepare",
			ID:            1,
			PrepareError:  sql.ErrConnDone,
			ExecError:     nil,
			ExpectedError: sql.ErrConnDone,
		},
		{
			Name:          "ErrorExec",
			ID:            1,
			PrepareError:  nil,
			ExecError:     sql.ErrConnDone,
			ExpectedError: sql.ErrConnDone,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.PrepareError != nil {
				mock.ExpectPrepare("UPDATE outbox SET status = 'done' WHERE id = \\$1").
					WillReturnError(test.PrepareError)
			} else {
				mock.ExpectPrepare("UPDATE outbox SET status = 'done' WHERE id = \\$1").
					ExpectExec().
					WithArgs(test.ID).
					WillReturnResult(sqlmock.NewResult(1, 1)).
					WillReturnError(test.ExecError)
			}

			err := service.SetDown(test.ID)

			if err != nil && (test.ExpectedError == nil) {
				t.Errorf("%s: expected error '%v', got '%v'", test.Name, test.ExpectedError, err)
			}

			if err == nil && test.ExpectedError != nil {
				t.Errorf("%s: expected error '%v', but got nil", test.Name, test.ExpectedError)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}

}

func TestMessageStorage_AddProcessedMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := storage.NewMessageStorage(db)

	mock.ExpectPrepare("UPDATE message SET status = 'update', processed_content = \\$1 WHERE id = \\$2").
		ExpectExec().
		WithArgs("Processed content", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.AddProcessedMessage(1, "Processed content")
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetDownMessages_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	service := storage.NewMessageStorage(db)

	rows := sqlmock.NewRows([]string{"id", "content", "processed_content", "status"}).
		AddRow(1, "test content", "processed content", "update").
		AddRow(2, "another content", "another processed content", "update")

	mock.ExpectQuery("SELECT id, content, processed_content, status FROM message WHERE status = 'update'").
		WillReturnRows(rows)

	expectedResult := []model.Message{
		{ID: 1, Content: "test content", ProcessedContent: "processed content", Status: "update"},
		{ID: 2, Content: "another content", ProcessedContent: "another processed content", Status: "update"},
	}

	result, err := service.GetDownMessages()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
