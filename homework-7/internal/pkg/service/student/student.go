package student

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/homework7/internal/pkg/model"
	"gitlab.ozon.dev/homework7/internal/pkg/repositories/studentsRepository"
)

func New(repo studentsRepository.StudentsRepo,
	cachedStudentsRepo studentsRepository.StudentsRepoCached) *Service {
	return &Service{
		repo:   repo,
		cached: cachedStudentsRepo,
	}
}

type Service struct {
	repo   studentsRepository.StudentsRepo
	cached studentsRepository.StudentsRepoCached
}

func (s *Service) Get(ctx context.Context, studentID uint64) (*model.Student, error) {
	student, err := s.cached.Get(ctx, studentID)
	if err == nil {
		return student, nil
	}
	if err != studentsRepository.ErrObjectNotFound {
		return nil, err
	}

	student, err = s.repo.GetById(ctx, studentID)
	if err != nil {
		fmt.Print(err)
	}
	return student, err
}

func (s *Service) Create(ctx context.Context, student *model.Student) (*model.Student, error) {
	dbStudent, err := s.repo.Add(ctx, student)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	s.cached.Add(ctx, student)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return dbStudent, nil
}

func (s *Service) Update(ctx context.Context, student *model.Student) error {
	isSuccess, err := s.repo.Update(ctx, student)
	if err != nil {
		fmt.Print(err)
		return err
	}
	if !isSuccess {
		fmt.Print("No student with this id")
		return studentsRepository.ErrObjectNotFound
	}

	s.cached.Add(ctx, student)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
