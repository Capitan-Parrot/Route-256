package grpcServer

import (
	"context"
	internal "gitlab.ozon.dev/homework8"
	"gitlab.ozon.dev/homework8/internal/pkg/model"
	studentService "gitlab.ozon.dev/homework8/internal/pkg/service/student"
	pb2 "gitlab.ozon.dev/homework8/pb"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Transport struct {
	pb2.UnimplementedStudentServiceServer
	studentService *studentService.Service
}

func NewImplementation(service *studentService.Service) *Transport {
	return &Transport{
		studentService: service,
	}
}

// CreateStudent Create создаёт профиль студента
func (t *Transport) CreateStudent(ctx context.Context, in *pb2.CreateStudentRequest) (*pb2.CreateStudentResponse, error) {
	tr := otel.Tracer("CreateStudent")
	ctx, span := tr.Start(ctx, "transport layer")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	student, err := t.studentService.CreateStudent(ctx, &model.Student{
		Name:          in.Name,
		CourseProgram: in.CourseProgram,
	})
	if err != nil {
		return nil, err
	}

	internal.CreateStudentCounter.Add(1)
	return &pb2.CreateStudentResponse{Id: student.ID}, nil
}

// GetStudent получает информацию о студенте по id
func (t *Transport) GetStudent(ctx context.Context, in *pb2.GetStudentRequest) (*pb2.GetStudentResponse, error) {
	tr := otel.Tracer("GetStudent")
	ctx, span := tr.Start(ctx, "transport layer")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	student, err := t.studentService.GetStudent(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	internal.GetStudentCounter.Add(1)
	return &pb2.GetStudentResponse{
		Id:            student.ID,
		Name:          student.Name,
		CourseProgram: student.CourseProgram,
		CreatedAt:     timestamppb.New(student.CreatedAt),
		UpdatedAt:     timestamppb.New(student.UpdatedAt),
	}, err
}

// UpdateStudent обновляет информацию о студенте
func (t *Transport) UpdateStudent(ctx context.Context, in *pb2.UpdateStudentRequest) (*pb2.UpdateStudentResponse, error) {
	tr := otel.Tracer("UpdateStudent")
	ctx, span := tr.Start(ctx, "transport layer")
	span.SetAttributes(attribute.Key("params").String(in.String()))
	defer span.End()

	ok, err := t.studentService.UpdateStudent(ctx, &model.Student{
		ID:            in.Id,
		Name:          in.Name,
		CourseProgram: in.CourseProgram,
	})

	return &pb2.UpdateStudentResponse{Ok: ok}, err
}
