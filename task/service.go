package task

import (
	"taskApp/db"

	"github.com/jmoiron/sqlx"
)

const (
	FindByPk = `
		SELECT task_pk, name, description, due_date, done FROM tasks WHERE task_pk = $1
		`
	FindAll = `
		SELECT task_pk, name, description, due_date, done FROM tasks 
		`
	Create = `
		INSERT INTO tasks (name, description, due_date, done) 
		VALUES (:name, :description, :due_date, :done)
		`
	Delete = `
		DELETE FROM tasks WHERE task_pk = $1
		`
	Update = `
		UPDATE tasks 
		SET name = :name, description = :description, due_date = :due_date, done = :done 
		WHERE task_pk = :task_pk
		`
)

type Service struct {
	sqlxDB *sqlx.DB
}

func NewService() *Service {
	database, _ := db.GetDB()
	return &Service{sqlxDB: sqlx.NewDb(database, "postgres")}
}

func (s *Service) getTask(taskPk int) (*Task, error) {
	task := &Task{}
	err := s.sqlxDB.Get(task, FindByPk, taskPk)
	return task, err
}

func (s *Service) getAllTasks() ([]*Task, error) {
	tasks := make([]*Task, 0)
	err := s.sqlxDB.Select(&tasks, FindAll)
	return tasks, err
}

func (s *Service) createTask(task *Task) error {
	_, err := s.sqlxDB.NamedExec(Create, task)
	return err
}

func (s *Service) deleteTask(taskPk int) error {
	_, err := s.sqlxDB.Exec(Delete, taskPk)
	return err
}

func (s *Service) updateTask(taskPk int, task *Task) error {
	task.TaskPk = taskPk
	_, err := s.sqlxDB.NamedExec(Update, task)
	return err
}
