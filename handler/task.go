package handler

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"my-todo/database"
	"my-todo/database/dbHelper"
	"my-todo/models"
	"my-todo/utils"
	"net/http"
)

//CreateTask creates a user's task using session token
func CreateTask(w http.ResponseWriter, r *http.Request) {

	body := models.CreateUserTask{}
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

	utils.RespondJSON(w, http.StatusCreated, models.CreateTaskRes{
		TaskID: taskID,
		Task:   createdTask,
	})
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("x-api-key")
	user, err := GetUserBySession(apiKey)
	if err != nil || user == nil {
		logrus.WithError(err).Errorf("failed to get user with token: %s", apiKey)
		utils.RespondError(w, http.StatusUnauthorized, err, "not authorized")
		return
	}
	var userTasks []models.UserTask
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		var saveErr error
		userTasks, saveErr = dbHelper.GetTask(user.ID)
		if saveErr != nil {
			return saveErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "error getting task")
		return
	}
	if userTasks == nil {
		userTasks = []models.UserTask{}
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Tasks []models.UserTask `json:"tasks"`
	}{
		userTasks,
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("x-api-key")
	user, err := GetUserBySession(apiKey)
	if err != nil || user == nil {
		logrus.WithError(err).Errorf("failed to get user with token: %s", apiKey)
		utils.RespondError(w, http.StatusUnauthorized, err, "not authorized")
		return
	}
	body := models.UpdateUserTask{}

	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
	}
	var UpdatedIsComplete bool
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		var saveErr error
		UpdatedIsComplete, saveErr = dbHelper.UpdateTask(tx, user.ID, body.TaskID, body.IsCompleted)
		if saveErr != nil {
			return saveErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "error updating task")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		IsCompleted bool `json:"isCompleted"`
	}{
		UpdatedIsComplete,
	})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("x-api-key")
	user, err := GetUserBySession(apiKey)
	if err != nil || user == nil {
		logrus.WithError(err).Errorf("failed to get user with token: %s", apiKey)
		utils.RespondError(w, http.StatusUnauthorized, err, "not authorized")
		return
	}
	body := models.UserTaskId{}

	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		saveErr := dbHelper.DeleteTask(tx, user.ID, body.TaskID)
		if saveErr != nil {
			return saveErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "error deleting task")
		return
	}
	utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{
		"success",
	})
}
