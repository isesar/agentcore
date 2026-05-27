package models

const (
	IntentGeneralQA      = "general_qa"
	IntentProjectLookup  = "project_lookup"
	IntentUserManagement = "user_management"
	IntentSmalltalk      = "smalltalk"
	IntentUnknown        = "unknown"

	ClassifierRules = "rules"
	ClassifierLLM   = "llm"

	RouterPolicyVersion = "v1"
)

type RouterDecision struct {
	Intent        string  `json:"intent" db:"intent"`
	Confidence    float64 `json:"confidence" db:"confidence"`
	Reason        string  `json:"reason" db:"reason"`
	Classifier    string  `json:"classifier" db:"classifier"`
	PolicyVersion string  `json:"policy_version" db:"policy_version"`
}
