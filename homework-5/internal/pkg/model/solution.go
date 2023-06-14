package model

import "time"

type Solution struct {
	ID        uint64    `db:"id" json:"id"`
	StudentID uint64    `db:"student_id" json:"studentId"`
	TaskID    uint64    `db:"task_id" json:"taskId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
