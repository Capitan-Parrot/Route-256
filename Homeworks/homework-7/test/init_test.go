package test

import (
	"github.com/bradfitz/gomemcache/memcache"
	"gitlab.ozon.dev/homework7/internal/config"
	cachedStudentsRepository "gitlab.ozon.dev/homework7/internal/pkg/repositories/studentsRepository/memcache"
	studentsRepository "gitlab.ozon.dev/homework7/internal/pkg/repositories/studentsRepository/postgresql"
	"gitlab.ozon.dev/homework7/internal/pkg/server"
	studentsService "gitlab.ozon.dev/homework7/internal/pkg/service/student"
	"gitlab.ozon.dev/homework7/test/postgres"
	"net/http/httptest"
)

var (
	Db          *postgres.TDB
	StudentRepo *studentsRepository.StudentsRepo
	Ts          *httptest.Server
)

func init() {
	mcClient := memcache.New(config.MemcachedPort)
	Db = postgres.NewFromEnv()
	// инициализация репозиториев + кэш
	StudentRepo = studentsRepository.NewStudents(Db.DB)
	cachedStudentRepo := cachedStudentsRepository.NewCachedRepo(mcClient)

	// сервисы для логики сервиса
	studentsServ := studentsService.New(StudentRepo, cachedStudentRepo)

	// транспорт для http запросов
	implementation := server.New(studentsServ)
	mux := implementation.RunServer()
	Ts = httptest.NewServer(mux)

}
