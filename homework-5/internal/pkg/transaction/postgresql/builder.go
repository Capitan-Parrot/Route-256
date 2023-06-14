package postgresql

import (
	"context"
	solutionsRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/solutionsRepository/postgresql"
	studentsRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/studentsRepository/postgresql"
	tasksRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/tasksRepository/postgresql"

	"gitlab.ozon.dev/homework5/internal/pkg/db"
	"gitlab.ozon.dev/homework5/internal/pkg/transaction"
)

type ServiceTxBuilder struct {
	db db.PGX
}

func NewServiceTxBuilder(db db.PGX) *ServiceTxBuilder {
	return &ServiceTxBuilder{db: db}
}

func (f *ServiceTxBuilder) ServiceTx(ctx context.Context) (*transaction.ServiceTx, error) {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &transaction.ServiceTx{
		Tx:       tx,
		Students: studentsRepository.NewStudents(tx),
		Tasks:    tasksRepository.NewTasks(tx),
		Solution: solutionsRepository.NewSolutions(tx),
	}, nil
}
