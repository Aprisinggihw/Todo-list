package router

import (
	"todo-list/internal/http/handler"
	"todo-list/pkg/route"
	"net/http"
)

func PublicRoutes(userHandler handler.UserHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: userHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: userHandler.Register,
		},
	}
}

func PrivateRoutes(userHandler handler.UserHandler, todosHandler handler.TodoHandler) []route.Route {
	return []route.Route{
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: userHandler.FindAll,
			Roles:   []string{"admin"},
		},
		{
			Method:  http.MethodPost,
			Path:    "/admin/todos",
			Handler: todosHandler.CreateTodoAsAdmin,
			Roles:   []string{"admin"},
		},
		{
			Method:  http.MethodPost,
			Path:    "/todos",
			Handler: todosHandler.CreateTodoHandler,
			Roles:   []string{"user"},
		},
		{
			Method:  http.MethodGet,
			Path:    "/admin/todos",
			Handler: todosHandler.GetAllHandler,
			Roles:   []string{"admin"},
		},
		{
			Method:  http.MethodPost,
			Path:    "/admin/todos",
			Handler: todosHandler.GetTodosByUserIdAsAdmin,
			Roles:   []string{"admin"},
		},
		{
			Method:  http.MethodGet,
			Path:    "/todos",
			Handler: todosHandler.GetTodosHandler,
			Roles:   []string{"user"},
		},
		{
			Method:  http.MethodPut,
			Path:    "/admin/user/:userID/todos/:todo_id",
			Handler: todosHandler.UpdateTodoAsAdmin,
			Roles:   []string{"admin"},
		},
		{
			Method:  http.MethodPut,
			Path:    "/todos/:id",
			Handler: todosHandler.UpdateTodoHandler,
			Roles:   []string{"user"},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/admin/user/:userID/todos/:todo_id",
			Handler: todosHandler.DeleteTodoAsAdmin,
			Roles:   []string{"admin"},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/todos/:id",
			Handler: todosHandler.DeleteTodoHandler,
			Roles:   []string{"user"},
		},
	}
}
