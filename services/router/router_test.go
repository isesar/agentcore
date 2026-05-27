package router

import (
	"errors"
	"testing"

	"agentcore/models"
)

type failingClassifier struct{}

func (f failingClassifier) Classify(prompt string, context models.QueryContextMetadata) (models.RouterDecision, error) {
	return models.RouterDecision{}, errors.New("classifier failed")
}

func TestRouteRulesHighConfidence(t *testing.T) {
	svc := NewService(false, 0.80, nil)
	decision, err := svc.Route(models.QueryRequest{Prompt: "show my project work"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Intent != models.IntentProjectLookup {
		t.Fatalf("expected %s, got %s", models.IntentProjectLookup, decision.Intent)
	}
	if decision.Classifier != models.ClassifierRules {
		t.Fatalf("expected rules classifier, got %s", decision.Classifier)
	}
	if decision.Reason == "" {
		t.Fatal("expected non-empty reason")
	}
}

func TestRouteLowConfidenceLLMDisabled(t *testing.T) {
	svc := NewService(false, 0.95, nil)
	decision, err := svc.Route(models.QueryRequest{Prompt: "what is the weather"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Intent != models.IntentUnknown {
		t.Fatalf("expected unknown intent, got %s", decision.Intent)
	}
	if decision.Reason != "rules_low_confidence_llm_disabled" {
		t.Fatalf("unexpected reason: %s", decision.Reason)
	}
}

func TestRouteLLMFallback(t *testing.T) {
	svc := NewService(true, 0.95, NewDeterministicLLMClassifier())
	decision, err := svc.Route(models.QueryRequest{Prompt: "what can you do"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if decision.Classifier != models.ClassifierLLM {
		t.Fatalf("expected llm classifier, got %s", decision.Classifier)
	}
	if decision.Intent == "" {
		t.Fatal("expected non-empty intent")
	}
}

func TestRouteLLMFailure(t *testing.T) {
	svc := NewService(true, 0.95, failingClassifier{})
	decision, err := svc.Route(models.QueryRequest{Prompt: "what can you do"})
	if err == nil {
		t.Fatal("expected error")
	}
	if decision.Intent != models.IntentUnknown {
		t.Fatalf("expected unknown intent, got %s", decision.Intent)
	}
	if decision.Reason != "llm_classifier_error" {
		t.Fatalf("unexpected reason: %s", decision.Reason)
	}
}
