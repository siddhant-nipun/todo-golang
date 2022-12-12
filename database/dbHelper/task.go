package dbHelper

import (
	"github.com/jmoiron/sqlx"
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

//func GetTask(db sqlx.Ext, userID string)(*models.UserTask, error){
//	//language=SQL
//	SQL:=`SELECT id, task, is_completed FROM users_task where user_id= $1`
//
//}
