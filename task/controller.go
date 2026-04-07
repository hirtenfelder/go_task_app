package task

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type ServiceInterface interface {
	getAllTasks() ([]*Task, error)
	getTask(taskPk int) (*Task, error)
	createTask(task *Task) error
	deleteTask(taskPk int) error
	updateTask(taskPk int, task *Task) error
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
		h.getAllTaskHandler(response, request)
	case http.MethodPost:
		h.createTaskHandler(response, request)
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
		h.getTaskHandler(response, request, taskPk)
	case http.MethodPut:
		h.updateTaskHandler(response, request, taskPk)
	case http.MethodDelete:
		h.deleteTaskHandler(response, request, taskPk)
	default:
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getAllTaskHandler(response http.ResponseWriter, _ *http.Request) {
	tasks, err := h.service.getAllTasks()
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

func (h *Handler) getTaskHandler(response http.ResponseWriter, _ *http.Request, taskPk int) {
	task, err := h.service.getTask(taskPk)
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

func (h *Handler) createTaskHandler(response http.ResponseWriter, request *http.Request) {
	task := &Task{}
	if json.NewDecoder(request.Body).Decode(task) != nil {
		h.handleEncoderDecoderError(response)
	}

	if err := h.service.createTask(task); err != nil {
		msg := "Failed to create task in database"
		slog.Warn(msg, "error", err)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteTaskHandler(response http.ResponseWriter, _ *http.Request, taskPk int) {
	if err := h.service.deleteTask(taskPk); err != nil {
		msg := "Failed to delete task from database"
		slog.Warn(msg, "error", err)
		http.Error(response, msg, http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateTaskHandler(response http.ResponseWriter, request *http.Request, taskPk int) {
	task := &Task{}
	if json.NewDecoder(request.Body).Decode(task) != nil {
		h.handleEncoderDecoderError(response)
	}

	if err := h.service.updateTask(taskPk, task); err != nil {
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
