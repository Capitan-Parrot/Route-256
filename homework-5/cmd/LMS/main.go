package main

import (
	"context"
	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.ozon.dev/homework5/internal/config"
	"gitlab.ozon.dev/homework5/internal/pkg/db"
	cachedSolutionsRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/solutionsRepository/memcache"
	solutionsRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/solutionsRepository/postgresql"
	cachedStudentsRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/studentsRepository/memcache"
	studentsRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/studentsRepository/postgresql"
	cachedTasksRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/tasksRepository/memcache"
	tasksRepository "gitlab.ozon.dev/homework5/internal/pkg/repositories/tasksRepository/postgresql"
	"gitlab.ozon.dev/homework5/internal/pkg/server"
	solutionsService "gitlab.ozon.dev/homework5/internal/pkg/service/solution"
	studentsService "gitlab.ozon.dev/homework5/internal/pkg/service/student"
	tasksService "gitlab.ozon.dev/homework5/internal/pkg/service/task"
	"gitlab.ozon.dev/homework5/internal/pkg/transaction/postgresql"
)

/*
Сервис, реализующий работу студента в LMS без прав админа.
Может регистрироваться, менять информацию о себе, получать список заданий,
отправлять решения и получать по ним текущий статус.
*/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		return
	}
	defer database.GetPool(ctx).Close()
	serviceTxBuilder := postgresql.NewServiceTxBuilder(database)
	mcClient := memcache.New(config.MemcachedPort)

	// инициализация репозиториев + кэш
	studentsRepo := studentsRepository.NewStudents(database)
	cachedStudentsRepo := cachedStudentsRepository.NewCachedRepo(mcClient)
	tasksRepo := tasksRepository.NewTasks(database)
	cachedTasksRepo := cachedTasksRepository.NewCachedRepo(mcClient)
	solutionsRepo := solutionsRepository.NewSolutions(database)
	cachedSolutionsRepo := cachedSolutionsRepository.NewCachedRepo(mcClient)

	// сервисы для логики сервиса
	studentsServ := studentsService.New(serviceTxBuilder, studentsRepo, cachedStudentsRepo)
	tasksServ := tasksService.New(serviceTxBuilder, tasksRepo, cachedTasksRepo)
	solutionsServ := solutionsService.New(serviceTxBuilder, solutionsRepo, cachedSolutionsRepo)

	// транспорт для http запросов
	implementation := server.New(studentsServ, tasksServ, solutionsServ)
	implementation.RunServer()
}
