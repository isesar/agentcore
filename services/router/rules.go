package router

import (
	"fmt"
	"regexp"
	"strings"
)

type rule struct {
	id      string
	intent  string
	pattern *regexp.Regexp
	weight  float64
}

func defaultRules() []rule {
	return []rule{
		{
			id:      "project_by_keyword",
			intent:  "project_lookup",
			pattern: regexp.MustCompile(`\b(project|projects|case study|work)\b`),
			weight:  0.90,
		},
		{
			id:      "user_management_by_keyword",
			intent:  "user_management",
			pattern: regexp.MustCompile(`\b(user|users|account|profile|delete user|create user)\b`),
			weight:  0.88,
		},
		{
			id:      "smalltalk_by_keyword",
			intent:  "smalltalk",
			pattern: regexp.MustCompile(`\b(hello|hi|hey|how are you|thanks|thank you)\b`),
			weight:  0.82,
		},
	}
}

func classifyByRules(prompt string, rules []rule) (string, float64, string) {
	lowerPrompt := strings.ToLower(prompt)
	for _, r := range rules {
		matches := r.pattern.FindAllString(lowerPrompt, -1)
		if len(matches) == 0 {
			continue
		}
		return r.intent, r.weight, fmt.Sprintf("rule=%s matched=%q", r.id, matches)
	}
	return "unknown", 0.30, "rule=no_match matched=[]"
}
