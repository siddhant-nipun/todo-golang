package models

type UserTask struct {
	ID          string `json:"id" db:"id"`
	Task        string `json:"task" db:"task"`
	IsCompleted string `json:"isCompleted" db:"is_completed"`
}
