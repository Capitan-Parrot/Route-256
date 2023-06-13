package tasksRepository

import (
	"context"
	"errors"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
)

var (
	ErrObjectNotFound = errors.New("task not found")
)

type TasksRepo interface {
	GetById(ctx context.Context, id uint64) (*model.Task, error)
	List(ctx context.Context) ([]*model.Task, error)
}

type TasksRepoCached interface {
	Add(ctx context.Context, student *model.Task) error
	Get(ctx context.Context, id uint64) (*model.Task, error)
	GetMulti(ctx context.Context, ids []*uint64) ([]*model.Task, error)
}
