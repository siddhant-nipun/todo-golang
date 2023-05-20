package dbHelper

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"my-todo/database"
	"my-todo/models"
	"my-todo/utils"
)

func IsUserExists(email string) (bool, error) {
	//language=SQL
	SQL := `SELECT id FROM users WHERE email = TRIM(LOWER($1))`
	var id string
	err := database.Todo.Get(&id, SQL, email)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, nil
	}
	return true, nil
}

func CreateUser(db sqlx.Ext, name, email, password string) (string, error) {
	//language=SQL
	SQL := `INSERT INTO users(name,email,password) VALUES ($1, TRIM(LOWER($2)), $3) RETURNING id`
	var userId string
	if err := db.QueryRowx(SQL, name, email, password).Scan(&userId); err != nil {
		return "", err
	}
	return userId, nil
}

func GetUserIDByPassword(email, password string) (string, error) {
	//	language=SQL
	SQL := `SELECT id, password FROM users WHERE email = TRIM(LOWER($1))`
	var userID string
	var passwordHash string

	err := database.Todo.QueryRowx(SQL, email).Scan(&userID, &passwordHash)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		return "", err
	}

	if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
		return "", errors.New("password error")
	}
	return userID, nil
}

func CreateUserSession(db sqlx.Ext, userID, sessionToken string) error {
	//language=SQL
	SQL := `INSERT INTO user_session(user_id, session_token) VALUES ($1,$2)`
	_, err := db.Exec(SQL, userID, sessionToken)
	return err
}

func DeleteUserSession(token string) error {
	//language=SQL
	SQL := `DELETE FROM user_session WHERE session_token=$1`
	_, err := database.Todo.Exec(SQL, token)
	return err
}

func GetUserBySession(sessionToken string) (*models.User, error) {
	//language="SQL"
	SQL := `SELECT u.id, u.name, u.email, u.created_at FROM users u 
			INNER JOIN user_session us on u.id = us.user_id
			WHERE us.session_token= $1`
	var user models.User
	err := database.Todo.Get(&user, SQL, sessionToken)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, nil
}
