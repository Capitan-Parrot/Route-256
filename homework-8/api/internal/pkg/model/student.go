package model

import (
	"time"
)

type Student struct {
	ID            uint64    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	CourseProgram string    `db:"course_program" json:"courseProgram"`
	CreatedAt     time.Time `db:"created_at" json:"-"`
	UpdatedAt     time.Time `db:"updated_at" json:"-"`
}
