package solutionsRepository

import (
	"context"
	"errors"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
)

var (
	ErrObjectNotFound = errors.New("solution not found")
)

type SolutionsRepo interface {
	GetById(ctx context.Context, id uint64) (*model.Solution, error)
	List(ctx context.Context) ([]*model.Solution, error)
	Add(ctx context.Context, student *model.Solution) (uint64, error)
}

type SolutionsRepoCached interface {
	Add(ctx context.Context, student *model.Solution) error
	Get(ctx context.Context, id uint64) (*model.Solution, error)
	GetMulti(ctx context.Context, ids []*uint64) ([]*model.Solution, error)
}
