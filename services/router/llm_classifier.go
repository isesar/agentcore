package router

import "agentcore/models"

type LLMClassifier interface {
	Classify(prompt string, context models.QueryContextMetadata) (models.RouterDecision, error)
}

type deterministicLLMClassifier struct{}

func NewDeterministicLLMClassifier() LLMClassifier {
	return &deterministicLLMClassifier{}
}

func (c *deterministicLLMClassifier) Classify(prompt string, _ models.QueryContextMetadata) (models.RouterDecision, error) {
	intent, confidence, _ := classifyByRules(prompt, defaultRules())
	if intent == models.IntentUnknown {
		return models.RouterDecision{
			Intent:        models.IntentGeneralQA,
			Confidence:    0.65,
			Reason:        "llm_fallback=heuristic_general_qa",
			Classifier:    models.ClassifierLLM,
			PolicyVersion: models.RouterPolicyVersion,
		}, nil
	}

	return models.RouterDecision{
		Intent:        intent,
		Confidence:    confidence,
		Reason:        "llm_fallback=heuristic_override",
		Classifier:    models.ClassifierLLM,
		PolicyVersion: models.RouterPolicyVersion,
	}, nil
}
