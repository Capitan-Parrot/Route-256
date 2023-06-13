package main

import (
	"context"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/homework8/config"
	"gitlab.ozon.dev/homework8/internal/pkg/db"
	"gitlab.ozon.dev/homework8/internal/pkg/grpcServer"
	cachedStudentsRepository "gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository/memcache"
	studentsRepository "gitlab.ozon.dev/homework8/internal/pkg/repositories/studentsRepository/postgresql"
	studentsService "gitlab.ozon.dev/homework8/internal/pkg/service/student"
	"gitlab.ozon.dev/homework8/pb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

/*
Сервис, реализующий работу студента в LMS без прав админа.
Может регистрироваться, получать и менять информацию о себе.
*/
func main() {
	tp, err := tracerProvider(config.TracesPort + "/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

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

	//// транспорт для http запросов
	//implementation := httpServer.New(studentsServ)
	//mux := implementation.RunServer()
	//if err := http.ListenAndServe(config.HttpPort, mux); err != nil {
	//	log.Fatal(err)
	//}

	// HTTP exporter для prometheus
	go http.ListenAndServe(config.PrometheusPort, promhttp.Handler())

	// транспорт для grpc запросов
	server := grpc.NewServer()
	pb.RegisterStudentServiceServer(server, grpcServer.NewImplementation(studentsServ))

	lsn, err := net.Listen("tcp", config.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting server on %s", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.Service),
			attribute.String("environment", config.Environment),
		)),
	)
	return tp, nil
}
