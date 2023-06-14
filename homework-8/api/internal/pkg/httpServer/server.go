package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.ozon.dev/homework8/internal/pkg/model"
	studentService "gitlab.ozon.dev/homework8/internal/pkg/service/student"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func New(studentService *studentService.Service) *Transport {
	return &Transport{
		studentService: studentService,
	}
}

type Transport struct {
	studentService *studentService.Service
}

func (t *Transport) RunServer() *http.ServeMux {
	mux := &http.ServeMux{}
	// эндпойнт для работы с таблицей Student
	mux.HandleFunc("/student", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			t.CreateStudent(res, req)
		case http.MethodGet:
			t.GetStudent(res, req)
		case http.MethodPut:
			t.UpdateStudent(res, req)
		default:
			t.Unsupported(res, req)
		}
	})

	return mux

}

// CreateStudent Create создаёт профиль студента
func (t *Transport) CreateStudent(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	student, err := getStudentFromBody(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbStudent, err := t.studentService.CreateStudent(ctx, student)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("Student.id: %d", dbStudent.ID)
}

// GetStudent получает информацию о студенте по id
func (t *Transport) GetStudent(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	studentID, err := getStudentID(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	student, err := t.studentService.GetStudent(ctx, studentID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonResult, err := json.Marshal(student)
	if err != nil {
		log.Println("error while marshalling student")
		return
	}
	fmt.Println(string(jsonResult))
}

// UpdateStudent обновляет информацию о студенте
func (t *Transport) UpdateStudent(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var student model.Student
	if err = json.Unmarshal(body, &student); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := t.studentService.UpdateStudent(ctx, &student); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("Success")
}

// Unsupported для неподдерживаемых методов
func (t *Transport) Unsupported(_ http.ResponseWriter, req *http.Request) {
	fmt.Printf("Unsupported method для server: %s", req.Method)
}

func getStudentID(url *url.URL) (uint64, error) {
	value := url.Query().Get(model.QueryParamStudentID)
	studentID, err := strconv.ParseUint(value, 10, 64)
	if err != nil || studentID == 0 {
		log.Println("invalid student.id")
		return 0, errors.New("invalid student.id")
	}
	return studentID, nil
}

func getStudentFromBody(body []byte) (*model.Student, error) {
	var student model.Student
	err := json.Unmarshal(body, &student)
	if err != nil || student.Name == "" {
		return nil, errors.New("invalid body")
	}
	return &student, nil
}
