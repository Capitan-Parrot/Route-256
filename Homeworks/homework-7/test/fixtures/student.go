package fixtures

import "gitlab.ozon.dev/homework7/internal/pkg/model"

type StudentBuilder struct {
	instance *model.Student
}

func Student() *StudentBuilder {
	return &StudentBuilder{instance: &model.Student{}}
}

func (s *StudentBuilder) ID(id uint64) *StudentBuilder {
	s.instance.ID = id
	return s
}

func (s *StudentBuilder) Name(name string) *StudentBuilder {
	s.instance.Name = name
	return s
}

func (s *StudentBuilder) CourseProgram(courseProgram string) *StudentBuilder {
	s.instance.CourseProgram = courseProgram
	return s
}

func (s *StudentBuilder) Valid() *StudentBuilder {
	return s.ID(1).Name("Andrew").CourseProgram("Go")
}

func (s *StudentBuilder) V() model.Student {
	return *s.instance
}

func (s *StudentBuilder) P() *model.Student {
	return s.instance
}
