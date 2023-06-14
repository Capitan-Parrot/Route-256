package student

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
	"gitlab.ozon.dev/homework6/internal/pkg/repositories/studentsRepository"
	"gitlab.ozon.dev/homework6/internal/pkg/transaction/postgresql"
)

func New(txModel *postgresql.ServiceTxBuilder, repo studentsRepository.StudentsRepo,
	cachedStudentsRepo studentsRepository.StudentsRepoCached) *Service {
	return &Service{
		txModel: txModel,
		repo:    repo,
		cached:  cachedStudentsRepo,
	}
}

type Service struct {
	txModel *postgresql.ServiceTxBuilder
	repo    studentsRepository.StudentsRepo
	cached  studentsRepository.StudentsRepoCached
}

func (s *Service) Get(ctx context.Context, studentID uint64) (*model.Student, error) {
	student, err := s.cached.Get(ctx, studentID)
	if err == nil {
		return student, nil
	}
	if err != studentsRepository.ErrObjectNotFound {
		return nil, err
	}

	tx, err := s.txModel.ServiceTx(ctx)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	student, err = tx.Students.GetById(ctx, studentID)
	if err != nil {
		fmt.Print(err)
	}
	return student, err
}

func (s *Service) Create(ctx context.Context, student *model.Student) (uint64, error) {
	tx, err := s.txModel.ServiceTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		fmt.Print(err)
		return 90, err
	}
	id, err := tx.Students.Add(ctx, student)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	tx.Commit(ctx)
	student.ID = id

	s.cached.Add(ctx, student)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	return id, nil
}

func (s *Service) Update(ctx context.Context, student *model.Student) error {
	tx, err := s.txModel.ServiceTx(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		fmt.Print(err)
		return err
	}
	isSuccess, err := tx.Students.Update(ctx, student)
	if err != nil {
		fmt.Print(err)
		return err
	}
	if !isSuccess {
		fmt.Print("No student with this id")
		return studentsRepository.ErrObjectNotFound
	}
	tx.Commit(ctx)

	s.cached.Add(ctx, student)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
