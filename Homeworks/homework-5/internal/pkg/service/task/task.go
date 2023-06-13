package task

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/homework5/internal/pkg/model"
	"gitlab.ozon.dev/homework5/internal/pkg/repositories/tasksRepository"
	"gitlab.ozon.dev/homework5/internal/pkg/transaction/postgresql"
)

func New(txModel *postgresql.ServiceTxBuilder, repo tasksRepository.TasksRepo,
	cachedTasksRepo tasksRepository.TasksRepoCached) *Service {
	return &Service{
		txModel: txModel,
		repo:    repo,
		cached:  cachedTasksRepo,
	}
}

type Service struct {
	txModel *postgresql.ServiceTxBuilder
	repo    tasksRepository.TasksRepo
	cached  tasksRepository.TasksRepoCached
}

func (s *Service) GetList(ctx context.Context) ([]*model.Task, error) {
	tx, err := s.txModel.ServiceTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	tasks, err := tx.Tasks.List(ctx)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	tx.Commit(ctx)
	return tasks, err
}
