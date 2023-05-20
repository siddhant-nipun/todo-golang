package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"my-todo/database"
	"my-todo/database/dbHelper"
	"my-todo/models"
	"my-todo/utils"
	"net/http"
	"time"
)

// RegisterUser to register a user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := models.RegisterRequest{}
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
	utils.RespondJSON(w, http.StatusCreated, models.TokenID{
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
	utils.RespondJSON(w, http.StatusOK, models.TokenID{
		Token: sessionToken,
	})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-api-key")
	user, err := dbHelper.GetUserBySession(token)
	if err != nil || user == nil {
		logrus.WithError(err).Errorf("failed to get user with token: %s", token)
		utils.RespondError(w, http.StatusUnauthorized, err, "not authorized")
		return
	}
	if err := dbHelper.DeleteUserSession(token); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout user")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
