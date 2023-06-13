//go:build integration

package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/homework8/test/fixtures"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	t.Run("success, db", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student := fixtures.Student().Valid().P()

		dbStudent, err := StudentRepo.Add(context.Background(), student)

		require.NoError(t, err)
		assert.Equal(t, dbStudent.Name, student.Name)
		assert.Equal(t, dbStudent.CourseProgram, student.CourseProgram)
	})

	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student := fixtures.Student().Name("more_than" + strings.Repeat("256", 256)).P()
		_, err := StudentRepo.Add(context.Background(), student)

		assert.Error(t, err)
	})
}

func TestGetById(t *testing.T) {
	t.Run("success, db", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student := fixtures.Student().Valid().P()
		dbStudent, err := StudentRepo.Add(context.Background(), student)

		fromDbStudent, err := StudentRepo.GetById(context.Background(), dbStudent.ID)

		require.NoError(t, err)
		assert.Equal(t, fromDbStudent, dbStudent)
	})

	t.Run("fail", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		dbStudent, err := StudentRepo.GetById(context.Background(), 0)

		require.Error(t, err)
		assert.Nil(t, dbStudent)
	})
}

func TestList(t *testing.T) {
	t.Run("success, db", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student1 := fixtures.Student().Valid().P()
		student2 := fixtures.Student().Valid().Name("Ivan").P()

		dbStudent1, err := StudentRepo.Add(context.Background(), student1)
		dbStudent2, err := StudentRepo.Add(context.Background(), student2)

		dbStudents, err := StudentRepo.List(context.Background())

		require.NoError(t, err)
		assert.Contains(t, dbStudents, dbStudent1)
		assert.Contains(t, dbStudents, dbStudent2)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success, db", func(t *testing.T) {
		Db.SetUp(t)
		defer Db.TearDown()

		student := fixtures.Student().Valid().P()
		newStudent, err := StudentRepo.Add(context.Background(), student)
		newStudent.CourseProgram = "C#"

		rowsAffected, err := StudentRepo.Update(context.Background(), newStudent)

		require.True(t, rowsAffected)
		dbStudent, err := StudentRepo.GetById(context.Background(), newStudent.ID)
		require.NoError(t, err)
		assert.Equal(t, dbStudent.CourseProgram, newStudent.CourseProgram)
	})
}
