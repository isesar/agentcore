package repository

import (
	"regexp"
	"testing"
	"time"

	"agentcore/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestMessageRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewMessageRepository(db)
	input := &models.Message{
		ConversationID: 1,
		Role:           "user",
		Content:        "hello",
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "conversation_id", "role", "content", "created_at"}).
		AddRow(10, input.ConversationID, input.Role, input.Content, now)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO messages (conversation_id, role, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, conversation_id, role, content, created_at`)).
		WithArgs(input.ConversationID, input.Role, input.Content).
		WillReturnRows(rows)

	created, err := repo.Create(input)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}
	if created.ID != 10 {
		t.Fatalf("expected message ID 10, got %d", created.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestMessageRepositoryListByConversationID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewMessageRepository(db)
	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "conversation_id", "role", "content", "created_at"}).
		AddRow(1, 77, "user", "hi", now).
		AddRow(2, 77, "assistant", "hello", now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, conversation_id, role, content, created_at
		 FROM messages
		 WHERE conversation_id = $1
		 ORDER BY id`)).
		WithArgs(77).
		WillReturnRows(rows)

	messages, err := repo.ListByConversationID(77)
	if err != nil {
		t.Fatalf("ListByConversationID() returned error: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[1].Role != "assistant" {
		t.Fatalf("unexpected second message role %q", messages[1].Role)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

var _ MessageRepository = (*PostgresMessageRepository)(nil)
