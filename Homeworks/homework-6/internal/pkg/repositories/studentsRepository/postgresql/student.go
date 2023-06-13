package postgresql

import (
	"context"
	"database/sql"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
	"gitlab.ozon.dev/homework6/internal/pkg/repositories/studentsRepository"
	"time"

	"gitlab.ozon.dev/homework6/internal/pkg/db"
)

type StudentsRepo struct {
	db db.DBops
}

func NewStudents(db db.DBops) *StudentsRepo {
	return &StudentsRepo{db: db}
}

// Add specific student
func (r *StudentsRepo) Add(ctx context.Context, student *model.Student) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO students(name, course_program) VALUES ($1, $2) RETURNING id`,
		student.Name, student.CourseProgram).Scan(&id)
	return id, err
}

func (r *StudentsRepo) GetById(ctx context.Context, id uint64) (*model.Student, error) {
	var u model.Student
	err := r.db.Get(ctx, &u, "SELECT id, name, course_program, created_at, updated_at FROM students WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return nil, studentsRepository.ErrObjectNotFound
	}
	return &u, err
}

func (r *StudentsRepo) List(ctx context.Context) ([]*model.Student, error) {
	students := make([]*model.Student, 0)
	err := r.db.Select(ctx, &students, "SELECT id,name, course_program, created_at,updated_at FROM students")
	return students, err
}

func (r *StudentsRepo) Update(ctx context.Context, student *model.Student) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE students SET name = $1, course_program = $2, updated_at = $3 WHERE id = $4",
		student.Name, student.CourseProgram, time.Now(), student.ID)
	return result.RowsAffected() > 0, err
}
