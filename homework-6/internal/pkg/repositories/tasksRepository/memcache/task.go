package memcache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/homework6/internal/config"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
	"gitlab.ozon.dev/homework6/internal/pkg/repositories/tasksRepository"
	"gitlab.ozon.dev/homework6/internal/pkg/stringUtils"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

type TasksRepo struct {
	cli *memcache.Client
}

func NewCachedRepo(cli *memcache.Client) *TasksRepo {
	return &TasksRepo{cli: cli}
}

func (t *TasksRepo) Add(_ context.Context, task *model.Task) error {
	jsonSolution, err := json.Marshal(task)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = t.cli.Set(&memcache.Item{
		Key:        strconv.FormatUint(task.ID, 10),
		Value:      jsonSolution,
		Expiration: config.TasksCacheExpiration,
	})
	return err
}

func (t *TasksRepo) Get(_ context.Context, id uint64) (*model.Task, error) {
	cachedTask, err := t.cli.Get(strconv.FormatUint(id, 10))
	if err == memcache.ErrCacheMiss {
		return nil, tasksRepository.ErrObjectNotFound
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var task model.Task
	err = json.Unmarshal(cachedTask.Value, &task)
	return &task, nil
}

func (t *TasksRepo) GetMulti(_ context.Context, ids []*uint64) ([]*model.Task, error) {
	stringsIDs := stringUtils.ConvertSlice(ids)
	it, err := t.cli.GetMulti(stringsIDs)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	tasks := make([]*model.Task, len(ids))
	index := 0
	for _, cachedTask := range it {
		var task model.Task
		err = json.Unmarshal(cachedTask.Value, &task)
		tasks[*ids[index]] = &task
	}
	return tasks, nil
}
