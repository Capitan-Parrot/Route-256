package memcache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/homework8/config"
	"gitlab.ozon.dev/homework8/internal/pkg/model"
	"gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

type StudentsRepo struct {
	cli *memcache.Client
}

func NewCachedRepo(cli *memcache.Client) *StudentsRepo {
	return &StudentsRepo{cli: cli}
}

func (s *StudentsRepo) Add(ctx context.Context, student *model.Student) error {
	tr := otel.Tracer("AddStudent")
	_, span := tr.Start(ctx, "cache layer")
	span.SetAttributes(
		attribute.Key("ID").Int64(int64(student.ID)),
		attribute.Key("Name").String(student.Name),
		attribute.Key("CourseProgram").String(student.CourseProgram),
		attribute.Key("CreatedAt").String(student.CreatedAt.String()),
		attribute.Key("UpdatedAt").String(student.UpdatedAt.String()),
	)
	defer span.End()

	jsonStudent, err := json.Marshal(student)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = s.cli.Set(&memcache.Item{
		Key:        strconv.FormatUint(student.ID, 10),
		Value:      jsonStudent,
		Expiration: config.StudentsCacheExpiration,
	})
	return err
}

func (s *StudentsRepo) Get(ctx context.Context, id uint64) (*model.Student, error) {
	tr := otel.Tracer("GetStudent")
	_, span := tr.Start(ctx, "cache layer")
	span.SetAttributes(attribute.Key("ID").Int64(int64(id)))
	defer span.End()

	it, err := s.cli.Get(strconv.FormatUint(id, 10))
	if err == memcache.ErrCacheMiss {
		return nil, studentsRepository.ErrObjectNotFound
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var student model.Student
	err = json.Unmarshal(it.Value, &student)
	return &student, err
}
