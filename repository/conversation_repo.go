package repository

import (
	"database/sql"

	"agentcore/models"
)

func (r *PostgresConversationRepository) Create(conversation *models.Conversation) (*models.Conversation, error) {
	var created models.Conversation
	err := r.db.QueryRow(
		`INSERT INTO conversations (user_id, title)
		 VALUES ($1, $2)
		 RETURNING id, user_id, title, created_at, updated_at`,
		conversation.UserID, conversation.Title,
	).Scan(&created.ID, &created.UserID, &created.Title, &created.CreatedAt, &created.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresConversationRepository) GetByID(id int) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.QueryRow(
		`SELECT id, user_id, title, created_at, updated_at
		 FROM conversations
		 WHERE id = $1`,
		id,
	).Scan(
		&conversation.ID,
		&conversation.UserID,
		&conversation.Title,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *PostgresConversationRepository) ListByUserID(userID int) ([]models.Conversation, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, title, created_at, updated_at
		 FROM conversations
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversations := make([]models.Conversation, 0)
	for rows.Next() {
		var conversation models.Conversation
		if err := rows.Scan(
			&conversation.ID,
			&conversation.UserID,
			&conversation.Title,
			&conversation.CreatedAt,
			&conversation.UpdatedAt,
		); err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}
