package services

import (
	"database/sql"

	"agentcore/db"
	"agentcore/models"
)

func GetAllUsers() ([]models.User, error) {
	rows, err := db.DB.Query(`SELECT id, name, email, created_at, updated_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	users := make([]models.User, 0)
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := db.DB.QueryRow(
		`SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(user *models.User) (*models.User, error) {
	var createdUser models.User
	err := db.DB.QueryRow(
		`INSERT INTO users (name, email) VALUES ($1, $2)
		 RETURNING id, name, email, created_at, updated_at`,
		user.Name, user.Email,
	).Scan(&createdUser.ID, &createdUser.Name, &createdUser.Email, &createdUser.CreatedAt, &createdUser.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func UpdateUser(id int, user *models.User) (*models.User, error) {
	var updatedUser models.User
	err := db.DB.QueryRow(
		`UPDATE users
		 SET name = $1, email = $2, updated_at = NOW()
		 WHERE id = $3
		 RETURNING id, name, email, created_at, updated_at`,
		user.Name, user.Email, id,
	).Scan(&updatedUser.ID, &updatedUser.Name, &updatedUser.Email, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func DeleteUser(id int) error {
	_, err := db.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	return err
}
