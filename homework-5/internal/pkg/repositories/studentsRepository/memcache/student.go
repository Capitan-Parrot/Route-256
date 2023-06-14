package memcache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/homework5/internal/config"
	"gitlab.ozon.dev/homework5/internal/pkg/model"
	"gitlab.ozon.dev/homework5/internal/pkg/repositories/studentsRepository"
	"gitlab.ozon.dev/homework5/internal/pkg/stringUtils"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

type StudentsRepo struct {
	cli *memcache.Client
}

func NewCachedRepo(cli *memcache.Client) *StudentsRepo {
	return &StudentsRepo{cli: cli}
}

func (s *StudentsRepo) Add(_ context.Context, student *model.Student) error {
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

func (s *StudentsRepo) Get(_ context.Context, id uint64) (*model.Student, error) {
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

func (s *StudentsRepo) GetMulti(_ context.Context, ids []*uint64) ([]*model.Student, error) {
	stringsIDs := stringUtils.ConvertSlice(ids)
	it, err := s.cli.GetMulti(stringsIDs)
	if err == memcache.ErrCacheMiss {
		return nil, studentsRepository.ErrObjectNotFound
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	students := make([]*model.Student, len(ids))
	index := 0
	for _, cachedStudent := range it {
		err = json.Unmarshal(cachedStudent.Value, students[*ids[index]])
	}
	return students, nil
}
