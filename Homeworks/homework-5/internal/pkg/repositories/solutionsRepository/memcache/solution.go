package memcache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/homework5/internal/config"
	"gitlab.ozon.dev/homework5/internal/pkg/model"
	"gitlab.ozon.dev/homework5/internal/pkg/repositories/solutionsRepository"
	"gitlab.ozon.dev/homework5/internal/pkg/stringUtils"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

type SolutionsRepo struct {
	cli *memcache.Client
}

func NewCachedRepo(cli *memcache.Client) *SolutionsRepo {
	return &SolutionsRepo{cli: cli}
}

func (s *SolutionsRepo) Add(_ context.Context, solution *model.Solution) error {
	jsonSolution, err := json.Marshal(solution)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = s.cli.Set(&memcache.Item{
		Key:        strconv.FormatUint(solution.ID, 10),
		Value:      jsonSolution,
		Expiration: config.SolutionsCacheExpiration,
	})
	return err
}

func (s *SolutionsRepo) Get(_ context.Context, id uint64) (*model.Solution, error) {
	it, err := s.cli.Get(strconv.FormatUint(id, 10))
	if err == memcache.ErrCacheMiss {
		return nil, solutionsRepository.ErrObjectNotFound
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var solution model.Solution
	err = json.Unmarshal(it.Value, &solution)
	return &solution, err

}

func (s *SolutionsRepo) GetMulti(_ context.Context, ids []*uint64) ([]*model.Solution, error) {
	stringsIDs := stringUtils.ConvertSlice(ids)
	it, err := s.cli.GetMulti(stringsIDs)
	if err == memcache.ErrCacheMiss {
		return nil, solutionsRepository.ErrObjectNotFound
	}
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	solutions := make([]*model.Solution, len(ids))
	index := 0
	for _, cachedSolution := range it {
		err = json.Unmarshal(cachedSolution.Value, solutions[*ids[index]])
	}
	return solutions, nil
}
