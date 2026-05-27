package repository

import "agentcore/models"

func (r *PostgresToolCallRepository) Create(toolCall *models.ToolCall) (*models.ToolCall, error) {
	var created models.ToolCall
	err := r.db.QueryRow(
		`INSERT INTO tool_calls (agent_run_id, name, arguments, result)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, agent_run_id, name, arguments, result, created_at`,
		toolCall.AgentRunID, toolCall.Name, toolCall.Arguments, toolCall.Result,
	).Scan(
		&created.ID,
		&created.AgentRunID,
		&created.Name,
		&created.Arguments,
		&created.Result,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresToolCallRepository) ListByAgentRunID(agentRunID int) ([]models.ToolCall, error) {
	rows, err := r.db.Query(
		`SELECT id, agent_run_id, name, arguments, result, created_at
		 FROM tool_calls
		 WHERE agent_run_id = $1
		 ORDER BY id`,
		agentRunID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	toolCalls := make([]models.ToolCall, 0)
	for rows.Next() {
		var toolCall models.ToolCall
		if err := rows.Scan(
			&toolCall.ID,
			&toolCall.AgentRunID,
			&toolCall.Name,
			&toolCall.Arguments,
			&toolCall.Result,
			&toolCall.CreatedAt,
		); err != nil {
			return nil, err
		}
		toolCalls = append(toolCalls, toolCall)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return toolCalls, nil
}
