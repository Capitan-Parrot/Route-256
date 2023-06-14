package model

import "time"

type Task struct {
	ID          uint64    `db:"id" json:"id"`
	Description string    `db:"description" json:"description"`
	Deadline    time.Time `db:"deadline" json:"deadline"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}
