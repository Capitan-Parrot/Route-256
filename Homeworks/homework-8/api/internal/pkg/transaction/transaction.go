package transaction

import (
	"context"
	"gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository"

	"gitlab.ozon.dev/homework8/internal/pkg/db"
)

type ServiceTxBuilder interface {
	ServiceTx(ctx context.Context) (*ServiceTx, error)
}

type ServiceTx struct {
	db.Tx
	Students studentsRepository.StudentsRepo
}
