package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"agentcore/models"
	"agentcore/repository"
	"github.com/gin-gonic/gin"
)

type stubRouteEngine struct {
	decision models.RouterDecision
	err      error
}

func (s stubRouteEngine) Route(req models.QueryRequest) (models.RouterDecision, error) {
	return s.decision, s.err
}

type stubRoutingDecisionRepo struct {
	createErr error
	created   []*models.RoutingDecisionRecord
}

func (s *stubRoutingDecisionRepo) Create(decision *models.RoutingDecisionRecord) (*models.RoutingDecisionRecord, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}
	s.created = append(s.created, decision)
	out := *decision
	out.ID = len(s.created)
	return &out, nil
}

func (s *stubRoutingDecisionRepo) ListByConversationID(conversationID string) ([]models.RoutingDecisionRecord, error) {
	return nil, nil
}

func (s *stubRoutingDecisionRepo) ListByIntent(intent string) ([]models.RoutingDecisionRecord, error) {
	return nil, nil
}

func (s *stubRoutingDecisionRepo) ListByTraceID(traceID string) ([]models.RoutingDecisionRecord, error) {
	return nil, nil
}

func TestQuerySuccess(t *testing.T) {
	originalRouterFactory := routerServiceFactory
	originalRepoFactory := routingDecisionRepositoryFactory
	defer func() {
		routerServiceFactory = originalRouterFactory
		routingDecisionRepositoryFactory = originalRepoFactory
	}()
	routerServiceFactory = func() routeEngine {
		return stubRouteEngine{
			decision: models.RouterDecision{
				Intent:        models.IntentProjectLookup,
				Confidence:    0.90,
				Reason:        "rule=project_by_keyword matched=[\"project\"]",
				Classifier:    models.ClassifierRules,
				PolicyVersion: models.RouterPolicyVersion,
			},
		}
	}
	repo := &stubRoutingDecisionRepo{}
	routingDecisionRepositoryFactory = func() repository.RoutingDecisionRepository { return repo }

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/query", Query)

	body := `{"prompt":"hello","context":{"conversation_id":"conv-1","trace_id":"trace-1","tags":["web"]}}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp models.QueryResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if resp.Answer == "" {
		t.Fatal("expected answer in response")
	}
	if resp.ConversationID != "conv-1" {
		t.Fatalf("expected conversation_id conv-1, got %q", resp.ConversationID)
	}
	if resp.Intent != models.IntentProjectLookup {
		t.Fatalf("expected intent %q, got %q", models.IntentProjectLookup, resp.Intent)
	}
	if resp.Confidence <= 0 {
		t.Fatal("expected positive confidence")
	}
	if resp.Reason == "" {
		t.Fatal("expected reason in response")
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected one persisted routing record, got %d", len(repo.created))
	}
}

func TestQueryMissingPromptReturnsEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/query", Query)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/query", strings.NewReader(`{"context":{"trace_id":"t1"}}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assertErrorEnvelope(t, w, http.StatusBadRequest, "validation_error")
}

func TestQueryMalformedJSONReturnsEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/query", Query)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/query", strings.NewReader(`{"prompt":"abc"`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assertErrorEnvelope(t, w, http.StatusBadRequest, "validation_error")
}

func TestQueryInternalFailureReturnsEnvelope(t *testing.T) {
	originalRouterFactory := routerServiceFactory
	original := queryResponder
	routerServiceFactory = func() routeEngine {
		return stubRouteEngine{
			decision: models.RouterDecision{
				Intent:        models.IntentGeneralQA,
				Confidence:    0.66,
				Reason:        "rule=test",
				Classifier:    models.ClassifierRules,
				PolicyVersion: models.RouterPolicyVersion,
			},
		}
	}
	queryResponder = func(req models.QueryRequest, decision models.RouterDecision) (models.QueryResponse, error) {
		return models.QueryResponse{}, errors.New("boom")
	}
	defer func() {
		routerServiceFactory = originalRouterFactory
		queryResponder = original
	}()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/query", Query)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/query", strings.NewReader(`{"prompt":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assertErrorEnvelope(t, w, http.StatusInternalServerError, "internal_error")
}

func TestQueryRoutingPersistenceFailureReturnsEnvelope(t *testing.T) {
	originalRouterFactory := routerServiceFactory
	originalRepoFactory := routingDecisionRepositoryFactory
	defer func() {
		routerServiceFactory = originalRouterFactory
		routingDecisionRepositoryFactory = originalRepoFactory
	}()
	routerServiceFactory = func() routeEngine {
		return stubRouteEngine{
			decision: models.RouterDecision{
				Intent:        models.IntentGeneralQA,
				Confidence:    0.70,
				Reason:        "rule=test",
				Classifier:    models.ClassifierRules,
				PolicyVersion: models.RouterPolicyVersion,
			},
		}
	}
	routingDecisionRepositoryFactory = func() repository.RoutingDecisionRepository {
		return &stubRoutingDecisionRepo{createErr: errors.New("db down")}
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/query", Query)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/query", strings.NewReader(`{"prompt":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assertErrorEnvelope(t, w, http.StatusInternalServerError, "internal_error")
}

func TestErrorEnvelopeConsistencyAcrossHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/users/:id", GetUser)
	router.POST("/api/v1/query", Query)

	tests := []struct {
		name         string
		method       string
		path         string
		body         string
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "user bad id",
			method:       http.MethodGet,
			path:         "/api/v1/users/not-a-number",
			expectedCode: http.StatusBadRequest,
			expectedErr:  "validation_error",
		},
		{
			name:         "query validation",
			method:       http.MethodPost,
			path:         "/api/v1/query",
			body:         `{"prompt":""}`,
			expectedCode: http.StatusBadRequest,
			expectedErr:  "validation_error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			if tc.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assertErrorEnvelope(t, w, tc.expectedCode, tc.expectedErr)
		})
	}
}

func assertErrorEnvelope(t *testing.T, w *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if w.Code != status {
		t.Fatalf("expected status %d, got %d", status, w.Code)
	}

	var envelope models.ErrorEnvelope
	if err := json.Unmarshal(w.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("failed to parse error envelope: %v", err)
	}
	if envelope.Error.Code != code {
		t.Fatalf("expected error code %q, got %q", code, envelope.Error.Code)
	}
	if envelope.Error.Message == "" {
		t.Fatal("expected non-empty error.message")
	}
}
