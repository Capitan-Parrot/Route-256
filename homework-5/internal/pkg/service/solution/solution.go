package solution

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/homework5/internal/pkg/model"
	"gitlab.ozon.dev/homework5/internal/pkg/repositories/solutionsRepository"
	"gitlab.ozon.dev/homework5/internal/pkg/transaction/postgresql"
)

func New(txModel *postgresql.ServiceTxBuilder, repo solutionsRepository.SolutionsRepo,
	cachedSolutionsRepo solutionsRepository.SolutionsRepoCached) *Service {
	return &Service{
		txModel: txModel,
		repo:    repo,
		cached:  cachedSolutionsRepo,
	}
}

type Service struct {
	txModel *postgresql.ServiceTxBuilder
	repo    solutionsRepository.SolutionsRepo
	cached  solutionsRepository.SolutionsRepoCached
}

func (s *Service) Get(ctx context.Context, solutionID uint64) (*model.Solution, error) {
	solution, err := s.cached.Get(ctx, solutionID)
	if err == nil {
		return solution, nil
	}
	if err != solutionsRepository.ErrObjectNotFound {
		return nil, err
	}

	tx, err := s.txModel.ServiceTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	solution, err = tx.Solution.GetById(ctx, solutionID)
	if err != nil {
		fmt.Print(err)
	}
	tx.Commit(ctx)
	return solution, err
}

func (s *Service) Create(ctx context.Context, solution *model.Solution) (uint64, error) {
	tx, err := s.txModel.ServiceTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	id, err := tx.Solution.Add(ctx, solution)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	tx.Commit(ctx)
	solution.ID = id

	s.cached.Add(ctx, solution)
	return id, nil
}
