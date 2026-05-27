package repository

import "agentcore/models"

func (r *PostgresUsageEventRepository) Create(event *models.UsageEvent) (*models.UsageEvent, error) {
	var created models.UsageEvent
	err := r.db.QueryRow(
		`INSERT INTO usage_events (agent_run_id, prompt_tokens, completion_tokens, total_tokens)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, agent_run_id, prompt_tokens, completion_tokens, total_tokens, created_at`,
		event.AgentRunID, event.PromptTokens, event.CompletionTokens, event.TotalTokens,
	).Scan(
		&created.ID,
		&created.AgentRunID,
		&created.PromptTokens,
		&created.CompletionTokens,
		&created.TotalTokens,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresUsageEventRepository) ListByAgentRunID(agentRunID int) ([]models.UsageEvent, error) {
	rows, err := r.db.Query(
		`SELECT id, agent_run_id, prompt_tokens, completion_tokens, total_tokens, created_at
		 FROM usage_events
		 WHERE agent_run_id = $1
		 ORDER BY id`,
		agentRunID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]models.UsageEvent, 0)
	for rows.Next() {
		var event models.UsageEvent
		if err := rows.Scan(
			&event.ID,
			&event.AgentRunID,
			&event.PromptTokens,
			&event.CompletionTokens,
			&event.TotalTokens,
			&event.CreatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
