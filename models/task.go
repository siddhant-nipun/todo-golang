package models

type UserTask struct {
	ID          string `json:"id" db:"id"`
	Task        string `json:"task" db:"task"`
	IsCompleted string `json:"isCompleted" db:"is_completed"`
}

type UpdateUserTask struct {
	TaskID      string `json:"taskId"`
	IsCompleted bool   `json:"isCompleted"`
}

type CreateUserTask struct {
	Task string `json:"task"`
}

type UserTaskId struct {
	TaskID string `json:"taskId"`
}

type CreateTaskRes struct {
	TaskID int    `json:"taskId"`
	Task   string `json:"task"`
}
