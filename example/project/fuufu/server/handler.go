package server

import (
	"context"

	"github.com/AugustineAurelius/fuufu/api"
)

var _ api.StrictServerInterface = &Handler{}

type Handler struct{}

//prints all methods in stdout
//go:generate go tool impl h*Handler api.StrictServerInterface

// Get Todo list
// (GET /api/v1/todo)
func (h *Handler) GetAllTodos(ctx context.Context, request api.GetAllTodosRequestObject) (_ api.GetAllTodosResponseObject, _ error) {
	panic("not implemented") // TODO: Implement
}

// Creates a new task
// (POST /api/v1/todo)
func (h *Handler) CreateNewTask(ctx context.Context, request api.CreateNewTaskRequestObject) (_ api.CreateNewTaskResponseObject, _ error) {
	panic("not implemented") // TODO: Implement
}

// delete task by todo_id
// (DELETE /api/v1/todo/{todo_id})
func (h *Handler) DeleteTaskByID(ctx context.Context, request api.DeleteTaskByIDRequestObject) (_ api.DeleteTaskByIDResponseObject, _ error) {
	panic("not implemented") // TODO: Implement
}

// Get task by todo_id
// (GET /api/v1/todo/{todo_id})
func (h *Handler) GetTaskByID(ctx context.Context, request api.GetTaskByIDRequestObject) (_ api.GetTaskByIDResponseObject, _ error) {
	panic("not implemented") // TODO: Implement
}
