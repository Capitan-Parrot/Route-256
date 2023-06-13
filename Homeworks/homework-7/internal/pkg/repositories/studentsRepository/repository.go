//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_stu
package studentsRepository

import (
	"context"
	"errors"
	"gitlab.ozon.dev/homework7/internal/pkg/model"
)

var (
	ErrObjectNotFound = errors.New("student not found")
)

type StudentsRepo interface {
	Add(ctx context.Context, student *model.Student) (*model.Student, error)
	GetById(ctx context.Context, id uint64) (*model.Student, error)
	List(ctx context.Context) ([]*model.Student, error)
	Update(ctx context.Context, student *model.Student) (bool, error)
}

type StudentsRepoCached interface {
	Add(ctx context.Context, student *model.Student) error
	Get(ctx context.Context, id uint64) (*model.Student, error)
	GetMulti(ctx context.Context, ids []*uint64) ([]*model.Student, error)
}
