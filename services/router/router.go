package router

import (
	"agentcore/models"
)

type Service struct {
	enableLLM  bool
	threshold  float64
	classifier LLMClassifier
	rules      []rule
}

func NewService(enableLLM bool, threshold float64, classifier LLMClassifier) *Service {
	if classifier == nil {
		classifier = NewDeterministicLLMClassifier()
	}
	return &Service{
		enableLLM:  enableLLM,
		threshold:  threshold,
		classifier: classifier,
		rules:      defaultRules(),
	}
}

func (s *Service) Route(req models.QueryRequest) (models.RouterDecision, error) {
	intent, confidence, reason := classifyByRules(req.Prompt, s.rules)
	rulesDecision := models.RouterDecision{
		Intent:        intent,
		Confidence:    confidence,
		Reason:        reason,
		Classifier:    models.ClassifierRules,
		PolicyVersion: models.RouterPolicyVersion,
	}

	if rulesDecision.Confidence >= s.threshold {
		return rulesDecision, nil
	}

	if s.enableLLM {
		llmDecision, err := s.classifier.Classify(req.Prompt, req.Context)
		if err != nil {
			return models.RouterDecision{
				Intent:        models.IntentUnknown,
				Confidence:    0.0,
				Reason:        "llm_classifier_error",
				Classifier:    models.ClassifierLLM,
				PolicyVersion: models.RouterPolicyVersion,
			}, err
		}
		return llmDecision, nil
	}

	return models.RouterDecision{
		Intent:        models.IntentUnknown,
		Confidence:    rulesDecision.Confidence,
		Reason:        "rules_low_confidence_llm_disabled",
		Classifier:    models.ClassifierRules,
		PolicyVersion: models.RouterPolicyVersion,
	}, nil
}
