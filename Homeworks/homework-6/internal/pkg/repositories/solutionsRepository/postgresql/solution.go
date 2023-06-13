package postgresql

import (
	"context"
	"database/sql"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
	"gitlab.ozon.dev/homework6/internal/pkg/repositories/solutionsRepository"

	"gitlab.ozon.dev/homework6/internal/pkg/db"
)

type SolutionsRepo struct {
	db db.DBops
}

func NewSolutions(db db.DBops) *SolutionsRepo {
	return &SolutionsRepo{db: db}
}

// Add specific solution
func (r *SolutionsRepo) Add(ctx context.Context, solution *model.Solution) (uint64, error) {
	var id uint64
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO solutions(student_id, task_id) VALUES ($1, $2) RETURNING id`,
		solution.StudentID, solution.TaskID).Scan(&id)
	return id, err
}

func (r *SolutionsRepo) GetById(ctx context.Context, id uint64) (*model.Solution, error) {
	var u model.Solution
	err := r.db.Get(ctx, &u,
		`SELECT id, student_id, task_id, status, created_at, updated_at 
FROM solutions WHERE id=$1`, id)
	if err == sql.ErrNoRows {
		return nil, solutionsRepository.ErrObjectNotFound
	}
	return &u, err
}

func (r *SolutionsRepo) List(ctx context.Context) ([]*model.Solution, error) {
	solutions := make([]*model.Solution, 0)
	err := r.db.Select(ctx, &solutions,
		"SELECT id, student_id, task_id, status, created_at, updated_at FROM solutions")
	return solutions, err
}
