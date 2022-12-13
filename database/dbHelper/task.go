package dbHelper

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"my-todo/database"
	"my-todo/models"
)

func CreateTask(db sqlx.Ext, userID, task string) (int, string, error) {
	//language=SQL
	SQL := `INSERT INTO users_task (user_id,task) VALUES ($1,$2) RETURNING id, task`
	var (
		rTask string
		rID   int
	)
	if err := db.QueryRowx(SQL, userID, task).Scan(&rID, &rTask); err != nil {
		return 0, "", err
	}
	return rID, rTask, nil
}

func GetTask(userID string) ([]models.UserTask, error) {
	//language=SQL
	SQL := `SELECT id, task, is_completed FROM users_task where user_id= $1`

	var userTasks []models.UserTask
	err := database.Todo.Select(&userTasks, SQL, userID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return userTasks, nil

}

func UpdateTask(db sqlx.Ext, userId, taskId string, isCompleted bool) (bool, error) {
	//language=SQL
	SQL := `UPDATE users_task SET is_completed = $1 where user_id=$2 AND id=$3 RETURNING is_completed`

	var updatedIsComplete bool
	err := db.QueryRowx(SQL, isCompleted, userId, taskId).Scan(&updatedIsComplete)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println(updatedIsComplete)
		return false, err
	}
	if err == sql.ErrNoRows {
		return false, err
	}
	return updatedIsComplete, nil
}

func DeleteTask(db sqlx.Ext, userId, taskId string) error {
	//language=SQL
	SQL := `DELETE FROM users_task WHERE user_id=$1 AND id= $2`

	_, err := db.Exec(SQL, userId, taskId)

	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == sql.ErrNoRows {
		return err
	}
	return nil
}
