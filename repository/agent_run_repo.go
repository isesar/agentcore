package repository

import (
	"database/sql"

	"agentcore/models"
)

func (r *PostgresAgentRunRepository) Create(run *models.AgentRun) (*models.AgentRun, error) {
	var created models.AgentRun
	err := r.db.QueryRow(
		`INSERT INTO agent_runs (conversation_id, model, status, started_at, finished_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, conversation_id, model, status, started_at, finished_at`,
		run.ConversationID, run.Model, run.Status, run.StartedAt, run.FinishedAt,
	).Scan(
		&created.ID,
		&created.ConversationID,
		&created.Model,
		&created.Status,
		&created.StartedAt,
		&created.FinishedAt,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresAgentRunRepository) GetByID(id int) (*models.AgentRun, error) {
	var run models.AgentRun
	err := r.db.QueryRow(
		`SELECT id, conversation_id, model, status, started_at, finished_at
		 FROM agent_runs
		 WHERE id = $1`,
		id,
	).Scan(&run.ID, &run.ConversationID, &run.Model, &run.Status, &run.StartedAt, &run.FinishedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &run, nil
}

func (r *PostgresAgentRunRepository) ListByConversationID(conversationID int) ([]models.AgentRun, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, model, status, started_at, finished_at
		 FROM agent_runs
		 WHERE conversation_id = $1
		 ORDER BY started_at DESC`,
		conversationID,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	runs := make([]models.AgentRun, 0)
	for rows.Next() {
		var run models.AgentRun
		if err := rows.Scan(&run.ID, &run.ConversationID, &run.Model, &run.Status, &run.StartedAt, &run.FinishedAt); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return runs, nil
}
