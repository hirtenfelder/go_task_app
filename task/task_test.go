package task

import "testing"

func TestTask_String(t *testing.T) {
	dueDate := "05.05.2026"
	task := Task{Name: "Test", DueDate: &dueDate}
	ident := task.String()
	expected := "Task{Name: Test, DueDate: 05.05.2026}"
	if ident != expected {
		t.Errorf("Expected %s, got %s", expected, ident)
	}
	t.Log(ident)
}
