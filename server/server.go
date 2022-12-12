package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"my-todo/handler"
	"my-todo/utils"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

//SetupRoutes to set up server routes
func SetupRoutes() *Server {
	router := chi.NewRouter()
	router.Route("/api", func(api chi.Router) {
		api.Get("/health", func(writer http.ResponseWriter, request *http.Request) {
			utils.RespondJSON(writer, http.StatusOK, struct {
				Status string `json:"status"`
			}{Status: "Hello golang"})
		})
		api.Route("/", func(public chi.Router) {
			public.Post("/register", handler.RegisterUser)
			public.Post("/login", handler.LoginUser)
		})
		api.Route("/task", func(task chi.Router) {
			task.Group(taskRoutes)
		})
	})

	return &Server{
		Router: router,
	}

}

//Run the server
func (srv *Server) Run(port string) error {
	srv.server = &http.Server{
		Addr:              port,
		Handler:           srv.Router,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	return srv.server.ListenAndServe()
}

//Shutdown the server
func (srv *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return srv.server.Shutdown(ctx)
}
