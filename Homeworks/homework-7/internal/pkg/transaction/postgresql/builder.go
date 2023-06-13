package postgresql

import (
	"context"
	"gitlab.ozon.dev/homework7/internal/pkg/db"
	studentsRepository "gitlab.ozon.dev/homework7/internal/pkg/repositories/studentsRepository/postgresql"
	"gitlab.ozon.dev/homework7/internal/pkg/transaction"
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
	}, nil
}
