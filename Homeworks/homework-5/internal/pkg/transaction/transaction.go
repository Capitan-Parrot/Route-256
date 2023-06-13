package transaction

import (
	"context"
	"gitlab.ozon.dev/homework5/internal/pkg/repositories/solutionsRepository"
	"gitlab.ozon.dev/homework5/internal/pkg/repositories/studentsRepository"
	"gitlab.ozon.dev/homework5/internal/pkg/repositories/tasksRepository"

	"gitlab.ozon.dev/homework5/internal/pkg/db"
)

type ServiceTxBuilder interface {
	ServiceTx(ctx context.Context) (*ServiceTx, error)
}

type ServiceTx struct {
	db.Tx
	Students studentsRepository.StudentsRepo
	Tasks    tasksRepository.TasksRepo
	Solution solutionsRepository.SolutionsRepo
}
