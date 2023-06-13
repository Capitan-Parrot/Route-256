package postgresql

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/homework8/internal/pkg/db"
	"gitlab.ozon.dev/homework8/internal/pkg/model"
	"gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type StudentsRepo struct {
	db db.DBops
}

func NewStudents(db db.DBops) *StudentsRepo {
	return &StudentsRepo{db: db}
}

// Add specific student
func (r *StudentsRepo) Add(ctx context.Context, student *model.Student) (*model.Student, error) {
	tr := otel.Tracer("AddStudent")
	ctx, span := tr.Start(ctx, "database layer")
	span.SetAttributes(
		attribute.Key("ID").Int64(int64(student.ID)),
		attribute.Key("Name").String(student.Name),
		attribute.Key("CourseProgram").String(student.CourseProgram),
		attribute.Key("CreatedAt").String(student.CreatedAt.String()),
		attribute.Key("UpdatedAt").String(student.UpdatedAt.String()),
	)
	defer span.End()

	var dbStudent model.Student
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO students(name, course_program, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name, course_program, created_at, updated_at`,
		student.Name, student.CourseProgram, student.CreatedAt, student.UpdatedAt).Scan(
		&dbStudent.ID, &dbStudent.Name, &dbStudent.CourseProgram,
		&dbStudent.CreatedAt, &dbStudent.UpdatedAt)
	return &dbStudent, err
}

func (r *StudentsRepo) GetById(ctx context.Context, id uint64) (*model.Student, error) {
	tr := otel.Tracer("GetStudent")
	ctx, span := tr.Start(ctx, "database layer")
	span.SetAttributes(attribute.Key("studentID").Int64(int64(id)))
	defer span.End()

	var u model.Student
	err := r.db.Get(ctx, &u, "SELECT id, name, course_program, created_at, updated_at FROM students WHERE id=$1", id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, studentsRepository.ErrObjectNotFound
	}
	return &u, err
}

func (r *StudentsRepo) List(ctx context.Context) ([]*model.Student, error) {
	tr := otel.Tracer("ListStudent")
	ctx, span := tr.Start(ctx, "database layer")
	defer span.End()

	students := make([]*model.Student, 0)
	err := r.db.Select(ctx, &students, "SELECT id,name, course_program, created_at,updated_at FROM students")
	return students, err
}

func (r *StudentsRepo) Update(ctx context.Context, student *model.Student) (bool, error) {
	tr := otel.Tracer("UpdateStudent")
	ctx, span := tr.Start(ctx, "database layer")
	span.SetAttributes(
		attribute.Key("ID").Int64(int64(student.ID)),
		attribute.Key("Name").String(student.Name),
		attribute.Key("CourseProgram").String(student.CourseProgram),
		attribute.Key("CreatedAt").String(student.CreatedAt.String()),
		attribute.Key("UpdatedAt").String(student.UpdatedAt.String()),
	)
	defer span.End()

	result, err := r.db.Exec(ctx,
		"UPDATE students SET name = $1, course_program = $2, updated_at = $3 WHERE id = $4",
		student.Name, student.CourseProgram, student.UpdatedAt, student.ID)
	return result.RowsAffected() > 0, err
}
