package main

import (
	"context"
	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.ozon.dev/homework7/internal/config"
	"gitlab.ozon.dev/homework7/internal/pkg/db"
	cachedStudentsRepository "gitlab.ozon.dev/homework7/internal/pkg/repositories/studentsRepository/memcache"
	studentsRepository "gitlab.ozon.dev/homework7/internal/pkg/repositories/studentsRepository/postgresql"
	"gitlab.ozon.dev/homework7/internal/pkg/server"
	studentsService "gitlab.ozon.dev/homework7/internal/pkg/service/student"
	"log"
	"net/http"
)

/*
Сервис, реализующий работу студента в LMS без прав админа.
Может регистрироваться, получать и менять информацию о себе.
*/
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	database, err := db.NewDB(ctx)
	if err != nil {
		return
	}
	defer database.GetPool(ctx).Close()
	mcClient := memcache.New(config.MemcachedPort)

	// инициализация репозиториев + кэш
	studentsRepo := studentsRepository.NewStudents(database)
	cachedStudentsRepo := cachedStudentsRepository.NewCachedRepo(mcClient)

	// сервисы для логики сервиса
	studentsServ := studentsService.New(studentsRepo, cachedStudentsRepo)

	// транспорт для http запросов
	implementation := server.New(studentsServ)
	mux := implementation.RunServer()
	if err := http.ListenAndServe(config.ServerPort, mux); err != nil {
		log.Fatal(err)
	}
}
