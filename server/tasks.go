package server

import (
	"github.com/go-chi/chi/v5"
	"my-todo/handler"
)

func taskRoutes(r chi.Router) {
	r.Group(func(task chi.Router) {
		task.Post("/", handler.CreateTask)
		//task.Get("/", handler.GetTasks)
	})
}
