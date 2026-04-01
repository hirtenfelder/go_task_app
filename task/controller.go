package task

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type ServiceInterface interface {
	GetAllTasks() ([]*Task, error)
	GetTask(taskPk int) (*Task, error)
	CreateTask(task *Task) error
	DeleteTask(taskPk int) error
	UpdateTask(taskPk int, task *Task) error
}

type Handler struct {
	service ServiceInterface
}

func init() {
	handler := NewHandler()
	http.HandleFunc("/tasks", handler.TaskHandler)
	http.HandleFunc("/tasks/{id}", handler.TaskOperationHandler)
}

func NewHandler() *Handler {
	return &Handler{service: NewService()}
}

func (h *Handler) TaskHandler(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		h.GetAllTaskHandler(response, request)
	case http.MethodPost:
		h.CreateTaskHandler(response, request)
	default:
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) TaskOperationHandler(response http.ResponseWriter, request *http.Request) {
	taskPkStr := request.PathValue("id")
	if taskPkStr == "" {
		msg := "Missing 'id' parameter"
		slog.Warn(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	taskPk, err := strconv.Atoi(taskPkStr)
	if err != nil {
		msg := "Invalid 'id' parameter"
		slog.Warn(msg)
		http.Error(response, msg, http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet:
		h.GetTaskHandler(response, request, taskPk)
	case http.MethodPut:
		h.UpdateTaskHandler(response, request, taskPk)
	case http.MethodDelete:
		h.DeleteTaskHandler(response, request, taskPk)
	default:
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) GetAllTaskHandler(response http.ResponseWriter, _ *http.Request) {
	tasks, err := h.service.GetAllTasks()
	if err != nil {
		msg := "Failed to retrieve tasks from database"
		slog.Warn(msg)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}

	// Return as JSON
	response.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(response).Encode(tasks) != nil {
		h.handleEncoderDecoderError(response)
	}
}

func (h *Handler) GetTaskHandler(response http.ResponseWriter, _ *http.Request, taskPk int) {
	task, err := h.service.GetTask(taskPk)
	if err != nil {
		msg := "Task not found"
		slog.Warn(msg)
		http.Error(response, msg, http.StatusNotFound)
		return
	}

	// Return as JSON
	response.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(response).Encode(task) != nil {
		h.handleEncoderDecoderError(response)
	}
}

func (h *Handler) CreateTaskHandler(response http.ResponseWriter, request *http.Request) {
	task := &Task{}
	if json.NewDecoder(request.Body).Decode(task) != nil {
		h.handleEncoderDecoderError(response)
	}

	if err := h.service.CreateTask(task); err != nil {
		msg := "Failed to create task in database"
		slog.Warn(msg, "error", err)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func (h *Handler) DeleteTaskHandler(response http.ResponseWriter, _ *http.Request, taskPk int) {
	if err := h.service.DeleteTask(taskPk); err != nil {
		msg := "Failed to delete task from database"
		slog.Warn(msg, "error", err)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateTaskHandler(response http.ResponseWriter, request *http.Request, taskPk int) {
	task := &Task{}
	if json.NewDecoder(request.Body).Decode(task) != nil {
		h.handleEncoderDecoderError(response)
	}

	if err := h.service.UpdateTask(taskPk, task); err != nil {
		msg := "Failed to update task in database"
		slog.Warn(msg, "error", err)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleEncoderDecoderError(response http.ResponseWriter) {
	msg := "Failed to encode/decode tasks as JSON"
	slog.Warn(msg)
	http.Error(response, msg, http.StatusInternalServerError)
}
