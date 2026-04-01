package task

import "fmt"

type Task struct {
	TaskPk      int     `db:"task_pk"`
	Name        string  `db:"name"`
	Description *string `db:"description"`
	DueDate     *string `db:"due_date"`
	Done        bool    `db:"done"`
}

func (l *Task) String() string {
	return fmt.Sprintf("Task{Name: %s, DueDate: %s}", l.Name, *l.DueDate)
}
