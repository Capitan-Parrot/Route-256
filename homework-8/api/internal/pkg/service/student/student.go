package student

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/homework8/internal/pkg/model"
	"gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"time"
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

func (s *Service) GetStudent(ctx context.Context, studentID uint64) (*model.Student, error) {
	tr := otel.Tracer("GetStudent")
	ctx, span := tr.Start(ctx, "service layer")
	span.SetAttributes(attribute.Key("studentID").Int64(int64(studentID)))
	defer span.End()

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

func (s *Service) CreateStudent(ctx context.Context, student *model.Student) (*model.Student, error) {
	tr := otel.Tracer("CreateStudent")
	ctx, span := tr.Start(ctx, "service layer")
	span.SetAttributes(
		attribute.Key("ID").Int64(int64(student.ID)),
		attribute.Key("Name").String(student.Name),
		attribute.Key("CourseProgram").String(student.CourseProgram),
		attribute.Key("CreatedAt").String(student.CreatedAt.String()),
		attribute.Key("UpdatedAt").String(student.UpdatedAt.String()),
	)
	defer span.End()

	student.CreatedAt = time.Now().UTC()
	student.UpdatedAt = time.Now().UTC()
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

func (s *Service) UpdateStudent(ctx context.Context, student *model.Student) (bool, error) {
	tr := otel.Tracer("UpdateStudent")
	ctx, span := tr.Start(ctx, "service layer")
	span.SetAttributes(
		attribute.Key("ID").Int64(int64(student.ID)),
		attribute.Key("Name").String(student.Name),
		attribute.Key("CourseProgram").String(student.CourseProgram),
		attribute.Key("CreatedAt").String(student.CreatedAt.String()),
		attribute.Key("UpdatedAt").String(student.UpdatedAt.String()),
	)
	defer span.End()

	student.UpdatedAt = time.Now().UTC()
	isSuccess, err := s.repo.Update(ctx, student)
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	if !isSuccess {
		fmt.Print("No student with this id")
		return false, studentsRepository.ErrObjectNotFound
	}

	s.cached.Add(ctx, student)
	if err != nil {
		fmt.Print(err)
		return false, err
	}
	return true, nil
}
