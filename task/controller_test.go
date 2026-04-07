package task

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockService struct{}

func (m *mockService) getAllTasks() ([]*Task, error) {
	dueDate := "05.05.2026"
	return []*Task{
		{TaskPk: 1, Name: "Test", DueDate: &dueDate},
	}, nil
}

func (m *mockService) getTask(taskPk int) (*Task, error) {
	dueDate := "05.05.2026"
	return &Task{TaskPk: taskPk, Name: "Test", DueDate: &dueDate}, nil
}

func (m *mockService) createTask(_ *Task) error {
	return nil
}
func (m *mockService) deleteTask(_ int) error {
	return nil
}
func (m *mockService) updateTask(_ int, _ *Task) error {
	return nil
}

func TestGetAllTaskHandler(t *testing.T) {
	handler := &Handler{service: &mockService{}}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()

	handler.getAllTaskHandler(rec, req)

	t.Run("StatusCode", func(t *testing.T) {
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("ContentType", func(t *testing.T) {
		contentType := rec.Header().Get("Content-Type")
		expected := "application/json"
		if contentType != expected {
			t.Errorf("Expected Content-Type %s, got %s", expected, contentType)
		}
	})

	t.Run("Body", func(t *testing.T) {
		body := strings.TrimSpace(rec.Body.String())
		expected := `[{"TaskPk":1,"Name":"Test","Description":null,"DueDate":"05.05.2026","Done":false}]`
		if body != expected {
			t.Errorf("Expected %s, got %s", expected, body)
		}
	})
}
