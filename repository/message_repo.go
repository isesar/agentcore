package repository

import (
	"database/sql"

	"agentcore/models"
)

func (r *PostgresMessageRepository) Create(message *models.Message) (*models.Message, error) {
	var created models.Message
	err := r.db.QueryRow(
		`INSERT INTO messages (conversation_id, role, content)
		 VALUES ($1, $2, $3)
		 RETURNING id, conversation_id, role, content, created_at`,
		message.ConversationID, message.Role, message.Content,
	).Scan(&created.ID, &created.ConversationID, &created.Role, &created.Content, &created.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PostgresMessageRepository) GetByID(id int) (*models.Message, error) {
	var message models.Message
	err := r.db.QueryRow(
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages
		 WHERE id = $1`,
		id,
	).Scan(&message.ID, &message.ConversationID, &message.Role, &message.Content, &message.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *PostgresMessageRepository) ListByConversationID(conversationID int) ([]models.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages
		 WHERE conversation_id = $1
		 ORDER BY id`,
		conversationID,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	messages := make([]models.Message, 0)
	for rows.Next() {
		var message models.Message
		if err := rows.Scan(&message.ID, &message.ConversationID, &message.Role, &message.Content, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
