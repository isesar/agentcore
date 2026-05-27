package repository

import (
	"regexp"
	"testing"
	"time"

	"agentcore/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestConversationRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewConversationRepository(db)
	userID := 1
	input := &models.Conversation{
		UserID: &userID,
		Title:  "Test conversation",
	}

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "created_at", "updated_at"}).
		AddRow(1, userID, input.Title, now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO conversations (user_id, title)
		 VALUES ($1, $2)
		 RETURNING id, user_id, title, created_at, updated_at`)).
		WithArgs(input.UserID, input.Title).
		WillReturnRows(rows)

	created, err := repo.Create(input)
	if err != nil {
		t.Fatalf("Create() returned error: %v", err)
	}

	if created.ID != 1 {
		t.Fatalf("expected conversation ID 1, got %d", created.ID)
	}
	if created.Title != input.Title {
		t.Fatalf("expected title %q, got %q", input.Title, created.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestConversationRepositoryGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	repo := NewConversationRepository(db)
	now := time.Now()
	userID := 1

	rows := sqlmock.NewRows([]string{"id", "user_id", "title", "created_at", "updated_at"}).
		AddRow(2, userID, "Existing conversation", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, title, created_at, updated_at
		 FROM conversations
		 WHERE id = $1`)).
		WithArgs(2).
		WillReturnRows(rows)

	conversation, err := repo.GetByID(2)
	if err != nil {
		t.Fatalf("GetByID() returned error: %v", err)
	}
	if conversation == nil {
		t.Fatal("expected conversation, got nil")
	}
	if conversation.ID != 2 {
		t.Fatalf("expected conversation ID 2, got %d", conversation.ID)
	}
	if conversation.Title != "Existing conversation" {
		t.Fatalf("unexpected title %q", conversation.Title)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

var _ ConversationRepository = (*PostgresConversationRepository)(nil)
