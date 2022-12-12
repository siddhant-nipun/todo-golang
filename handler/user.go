package handler

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"my-todo/database"
	"my-todo/database/dbHelper"
	"my-todo/models"
	"my-todo/utils"
	"net/http"
	"time"
)

// RegisterUser to register a user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}
	exists, existsErr := dbHelper.IsUserExists(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}
	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
	}

	sessionToken := utils.HashString(body.Email + time.Now().String())

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
		if saveErr != nil {
			return saveErr
		}
		sessionErr := dbHelper.CreateUserSession(tx, userID, sessionToken)
		if sessionErr != nil {
			return sessionErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
}

// LoginUser to log in the user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	body := models.UserCredentials{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	userID, userErr := dbHelper.GetUserIDByPassword(body.Email, body.Password)
	if userErr != nil {
		utils.RespondError(w, http.StatusBadRequest, userErr, "login failed")
		return
	}

	sessionToken := utils.HashString(body.Email + time.Now().String())
	if sessionErr := dbHelper.CreateUserSession(database.Todo, userID, sessionToken); sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, nil, "failed to create user session")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
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
