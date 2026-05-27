package repository

import "agentcore/models"

func (r *PostgresErrorEventRepository) Create(event *models.ErrorEvent) (*models.ErrorEvent, error) {
	var created models.ErrorEvent
	err := r.db.QueryRow(
		`INSERT INTO errors (agent_run_id, error_type, message, stack_trace)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, agent_run_id, error_type, message, stack_trace, created_at`,
		event.AgentRunID, event.ErrorType, event.Message, event.StackTrace,
	).Scan(
		&created.ID,
		&created.AgentRunID,
		&created.ErrorType,
		&created.Message,
		&created.StackTrace,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresErrorEventRepository) ListByAgentRunID(agentRunID int) ([]models.ErrorEvent, error) {
	rows, err := r.db.Query(
		`SELECT id, agent_run_id, error_type, message, stack_trace, created_at
		 FROM errors
		 WHERE agent_run_id = $1
		 ORDER BY id`,
		agentRunID,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	events := make([]models.ErrorEvent, 0)
	for rows.Next() {
		var event models.ErrorEvent
		if err := rows.Scan(
			&event.ID,
			&event.AgentRunID,
			&event.ErrorType,
			&event.Message,
			&event.StackTrace,
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
