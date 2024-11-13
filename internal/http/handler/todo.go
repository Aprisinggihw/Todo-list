package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"todo-list/internal/entity"
	"todo-list/internal/service"
	"todo-list/pkg/response"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService}
}

func (h *TodoHandler) CreateTodoAsAdmin(ctx echo.Context) error {
	var req struct {
		UserID uint `json:"user_id"`
		Title string `json:"title"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	todo, err := h.todoService.CreateTodo(ctx.Request().Context(), req.UserID, req.Title)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("todo created successfully", todo))
}

func (h *TodoHandler) CreateTodoHandler(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized,response.ErrorResponse( http.StatusUnauthorized,fmt.Sprintf("Invalid or missing userID: %d " ,userID) ))
	}
	var req struct {
		Title string `json:"title"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	todo, err := h.todoService.CreateTodo(ctx.Request().Context(), userID, req.Title)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("todo created successfully", todo))
}

func (h *TodoHandler) GetAllHandler(ctx echo.Context) error {

	todos, err := h.todoService.GetTodos(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("successfully fetch all todos", todos))

}

func (h *TodoHandler) GetTodosByUserIdAsAdmin(ctx echo.Context) error {
	var reqUserID uint
	if err := ctx.Bind(reqUserID); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	
	todos, err := h.todoService.GetTodosByUserID(ctx.Request().Context(), reqUserID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("successfully fetch all todos", todos))

}

func (h *TodoHandler) GetTodosHandler(ctx echo.Context) error {

	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized,response.ErrorResponse( http.StatusUnauthorized,fmt.Sprintf("Invalid or missing userID: %d " ,userID) ))
	}
	todos, err := h.todoService.GetTodosByUserID(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("successfully fetch all todos", todos))

}

func (h *TodoHandler) UpdateTodoAsAdmin(ctx echo.Context) error {

	todoID, err := strconv.ParseUint(ctx.Param("todo_id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid todo ID"))
	}

	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid todo ID"))
	}

	req := new(entity.Todo)
	req.ID = uint(todoID)
	req.UserID = uint(userID)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err = h.todoService.UpdateTodo(context.Background(), req.UserID, req.ID, req.Title, req.Done)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("todo updated successfully", req))
}

func (h *TodoHandler) UpdateTodoHandler(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "Invalid or missing userID"))
	}

	todoID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid todo ID"))
	}

	var req struct {
		Title string `json:"title"`
		Done  bool   `json:"done"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	err = h.todoService.UpdateTodo(context.Background(), userID, uint(todoID), req.Title, req.Done)
	if err != nil {
		if err.Error() == "unauthorized or not found" {
			return ctx.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("todo updated successfully", req))
}

func (h *TodoHandler) DeleteTodoAsAdmin(ctx echo.Context) error {
	
	todoID, err := strconv.ParseUint(ctx.Param("todo_id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid todo ID"))
	}
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid todo ID"))
	}

	err = h.todoService.DeleteTodo(ctx.Request().Context(), uint(userID), uint(todoID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("todo deleted successfully", nil))
}


func (h *TodoHandler) DeleteTodoHandler(ctx echo.Context) error {
	userID, ok := ctx.Get("user_id").(uint)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized,response.ErrorResponse( http.StatusUnauthorized,fmt.Sprintf("Invalid or missing userID: %d " ,userID) ))
	}
	todoID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, "Invalid todo ID"))
	}
	var req struct {
		TodoID uint 
	}
	req.TodoID = uint(todoID)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}
	err = h.todoService.DeleteTodo(ctx.Request().Context(), userID, req.TodoID)
	if err != nil {
		if err.Error() == "unauthorized or not found" {
			return ctx.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, err.Error()))
		}
		return ctx.JSON(http.StatusInternalServerError, response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("todo deleted successfully", nil))
}
