package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"my-todo/database"
	"my-todo/database/dbHelper"
	"my-todo/utils"
	"net/http"
)

//CreateTask creates a user's task using session token
func CreateTask(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Task string `json:"task"`
	}{}
	apiKey := r.Header.Get("x-api-key")
	user, err := GetUserBySession(apiKey)
	if err != nil || user == nil {
		logrus.WithError(err).Errorf("failed to get user with token: %s", apiKey)
		utils.RespondError(w, http.StatusUnauthorized, err, "not authorized")
		return
	}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if err != nil || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, err, "invalid token")
		return
	}
	var (
		taskID      int
		createdTask string
	)
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		var saveErr error
		taskID, createdTask, saveErr = dbHelper.CreateTask(tx, user.ID, body.Task)
		if saveErr != nil {
			return saveErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "error creating task")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, struct {
		TaskID int    `json:"taskId"`
		Task   string `json:"task"`
	}{
		TaskID: taskID,
		Task:   createdTask,
	})
}

//func GetTasks(w http.ResponseWriter, r *http.Request) {
//	apiKey := r.Header.Get("x-api-key")
//	user, err := GetUserBySession(apiKey)
//	if err != nil || user == nil {
//		logrus.WithError(err).Errorf("failed to get user with token: %s", apiKey)
//		utils.RespondError(w, http.StatusUnauthorized, err, "not authorized")
//		return
//	}
//	//var UserTask = make(models.UserTask)
//	txErr := database.Tx(func(tx *sqlx.Tx) error {
//		var saveErr error
//		, saveErr = dbHelper.GetTask(tx, user.ID, body.Task)
//		if saveErr != nil {
//			return saveErr
//		}
//		return nil
//	})
//	if txErr != nil {
//		utils.RespondError(w, http.StatusInternalServerError, txErr, "error creating task")
//		return
//	}
//
//}
