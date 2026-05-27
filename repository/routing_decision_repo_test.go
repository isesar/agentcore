package repository

import (
	"encoding/json"
	"regexp"
	"testing"
	"time"

	"agentcore/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestRoutingDecisionRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = db.Close() }()

	repo := NewRoutingDecisionRepository(db)
	record := &models.RoutingDecisionRecord{
		ConversationID: "conv-1",
		Intent:         models.IntentProjectLookup,
		Confidence:     0.90,
		Reason:         "rule=project_by_keyword matched=[\"project\"]",
		Classifier:     models.ClassifierRules,
		PolicyVersion:  models.RouterPolicyVersion,
		TraceID:        "trace-1",
		UserID:         "user-1",
		Source:         "web",
		Tags:           []string{"agentcore", "query"},
	}

	tagsJSON, _ := json.Marshal(record.Tags)
	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "conversation_id", "agent_run_id", "intent", "confidence", "reason", "classifier",
		"policy_version", "trace_id", "user_id", "source", "tags", "created_at",
	}).AddRow(
		1, record.ConversationID, nil, record.Intent, record.Confidence, record.Reason, record.Classifier,
		record.PolicyVersion, record.TraceID, record.UserID, record.Source, tagsJSON, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO routing_decisions (
			conversation_id, agent_run_id, intent, confidence, reason, classifier,
			policy_version, trace_id, user_id, source, tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
		RETURNING id, conversation_id, agent_run_id, intent, confidence, reason, classifier, policy_version, trace_id, user_id, source, tags, created_at`)).
		WithArgs(record.ConversationID, record.AgentRunID, record.Intent, record.Confidence, record.Reason, record.Classifier, record.PolicyVersion, record.TraceID, record.UserID, record.Source, string(tagsJSON)).
		WillReturnRows(rows)

	created, err := repo.Create(record)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != 1 {
		t.Fatalf("expected ID 1, got %d", created.ID)
	}
	if len(created.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(created.Tags))
	}
}

func TestRoutingDecisionRepositoryListByTraceID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer func() { _ = db.Close() }()

	repo := NewRoutingDecisionRepository(db)
	now := time.Now()
	tagsJSON, _ := json.Marshal([]string{"agentcore"})
	rows := sqlmock.NewRows([]string{
		"id", "conversation_id", "agent_run_id", "intent", "confidence", "reason", "classifier",
		"policy_version", "trace_id", "user_id", "source", "tags", "created_at",
	}).AddRow(
		1, "conv-1", nil, models.IntentProjectLookup, 0.9, "rule=test", models.ClassifierRules,
		models.RouterPolicyVersion, "trace-1", "user-1", "web", tagsJSON, now,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, conversation_id, agent_run_id, intent, confidence, reason, classifier, policy_version, trace_id, user_id, source, tags, created_at
		 FROM routing_decisions
		 WHERE trace_id = $1
		 ORDER BY created_at DESC`)).
		WithArgs("trace-1").
		WillReturnRows(rows)

	decisions, err := repo.ListByTraceID("trace-1")
	if err != nil {
		t.Fatalf("ListByTraceID returned error: %v", err)
	}
	if len(decisions) != 1 {
		t.Fatalf("expected 1 decision, got %d", len(decisions))
	}
	if decisions[0].Intent != models.IntentProjectLookup {
		t.Fatalf("unexpected intent: %s", decisions[0].Intent)
	}
}
