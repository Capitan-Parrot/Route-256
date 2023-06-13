package postgresql

import (
	"context"
	"database/sql"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
	"gitlab.ozon.dev/homework6/internal/pkg/repositories/studentsRepository"

	"gitlab.ozon.dev/homework6/internal/pkg/db"
)

type TasksRepo struct {
	db db.DBops
}

func NewTasks(db db.DBops) *TasksRepo {
	return &TasksRepo{db: db}
}

func (r *TasksRepo) GetById(ctx context.Context, id uint64) (*model.Task, error) {
	var u model.Task
	err := r.db.Get(ctx, &u,
		"SELECT id, description, deadline, created_at, updated_at FROM tasks WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, studentsRepository.ErrObjectNotFound
	}
	return &u, err
}

func (r *TasksRepo) List(ctx context.Context) ([]*model.Task, error) {
	tasks := make([]*model.Task, 0)
	err := r.db.Select(ctx, &tasks,
		"SELECT id, description, deadline, created_at, updated_at FROM tasks")
	return tasks, err
}
