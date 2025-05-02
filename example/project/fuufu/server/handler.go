package server

import (
	"context"

	"github.com/AugustineAurelius/fuufu/api"
)

var _ api.StrictServerInterface = &Handler{}

//go:generate go run github.com/AugustineAurelius/eos/ generator wrapper  --name Handler
type Handler struct{}

//prints all methods in stdout
//go:generate go tool impl h*Handler api.StrictServerInterface

// Get Todo list
// (GET /api/v1/todo)
func (h *Handler) GetAllTodos(ctx context.Context, request api.GetAllTodosRequestObject) (api.GetAllTodosResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Creates a new task
// (POST /api/v1/todo)
func (h *Handler) CreateNewTask(ctx context.Context, request api.CreateNewTaskRequestObject) (api.CreateNewTaskResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// delete task by todo_id
// (DELETE /api/v1/todo/{todo_id})
func (h *Handler) DeleteTaskByID(ctx context.Context, request api.DeleteTaskByIDRequestObject) (api.DeleteTaskByIDResponseObject, error) {
	panic("not implemented") // TODO: Implement
}

// Get task by todo_id
// (GET /api/v1/todo/{todo_id})
func (h *Handler) GetTaskByID(ctx context.Context, request api.GetTaskByIDRequestObject) (api.GetTaskByIDResponseObject, error) {
	panic("not implemented") // TODO: Implement
}
