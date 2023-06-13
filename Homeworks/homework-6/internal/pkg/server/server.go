package server

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.ozon.dev/homework6/internal/config"
	"gitlab.ozon.dev/homework6/internal/pkg/model"
	consoleService "gitlab.ozon.dev/homework6/internal/pkg/service/console"
	solutionService "gitlab.ozon.dev/homework6/internal/pkg/service/solution"
	studentService "gitlab.ozon.dev/homework6/internal/pkg/service/student"
	taskService "gitlab.ozon.dev/homework6/internal/pkg/service/task"
	"io"
	"log"
	"net/http"
	"strconv"
)

func New(consoleService *consoleService.Service, studentService *studentService.Service,
	taskService *taskService.Service, solutionService *solutionService.Service) *Transport {
	return &Transport{
		consoleService:  consoleService,
		studentService:  studentService,
		tasksService:    taskService,
		solutionService: solutionService,
	}
}

type Transport struct {
	consoleService  *consoleService.Service
	studentService  *studentService.Service
	tasksService    *taskService.Service
	solutionService *solutionService.Service
}

func (t *Transport) RunServer() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/help", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			t.Help(res, req)
		default:
			t.Unsupported(res, req)
		}
	})
	mux.HandleFunc("/spell", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			t.Spell(res, req)
		default:
			t.Unsupported(res, req)
		}
	})
	mux.HandleFunc("/goFmt", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			t.GoFmt(res, req)
		default:
			t.Unsupported(res, req)
		}
	})
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

	// эндпойнт для работы с таблицей Task
	mux.HandleFunc("/task", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			t.GetTask(res, req)
		default:
			t.Unsupported(res, req)
		}
	})

	// эндпойнт для работы с таблицей Solution
	mux.HandleFunc("/solution", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			t.GetSolution(res, req)
		case http.MethodPost:
			t.CreateSolution(res, req)
		default:
			t.Unsupported(res, req)
		}
	})
	if err := http.ListenAndServe(config.ServerPort, mux); err != nil {
		log.Fatal(err)
	}
}

func (t *Transport) Help(w http.ResponseWriter, req *http.Request) {
	t.consoleService.Help()
}

func (t *Transport) Spell(w http.ResponseWriter, req *http.Request) {
	word := req.URL.Query().Get(model.QueryParamSpellKey)
	err := t.consoleService.Spell(word)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (t *Transport) GoFmt(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(config.GoFmtFileMemoryMax)

	file, _, err := req.FormFile(model.QueryParamGoFmtKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	bytes, err := t.consoleService.GoFmt(file)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(bytes)
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

	var student model.Student
	if err = json.Unmarshal(body, &student); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := t.studentService.Create(ctx, &student)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("Student.id: %d", id)
}

// GetStudent получает информацию о студенте по id
func (t *Transport) GetStudent(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	value := req.URL.Query().Get(model.QueryParamStudentID)
	studentID, err := strconv.ParseUint(value, 10, 64)
	if err != nil || studentID == 0 {
		log.Println("invalid student.id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	student, err := t.studentService.Get(ctx, studentID)
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

	if err = t.studentService.Update(ctx, &student); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("Success")
}

// GetTask Список текущих задач
func (t *Transport) GetTask(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasks, err := t.tasksService.GetList(ctx)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonResult, err := json.Marshal(tasks)
	if err != nil {
		log.Println("error while marshalling")
		return
	}
	fmt.Println(string(jsonResult))
}

// CreateSolution посылает решение студента
func (t *Transport) CreateSolution(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var solution model.Solution
	if err = json.Unmarshal(body, &solution); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := t.solutionService.Create(ctx, &solution)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("Solution.id: %d", id)
}

// GetSolution отдаёт статус по решению
func (t *Transport) GetSolution(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	value := req.URL.Query().Get(model.QueryParamSolutionID)
	solutionID, err := strconv.ParseUint(value, 10, 64)
	if err != nil || solutionID == 0 {
		log.Println("invalid solution.id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err := t.solutionService.Get(ctx, solutionID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonResult, err := json.Marshal(task)
	if err != nil {
		log.Println("error while marshalling task")
		return
	}
	fmt.Println(string(jsonResult))
}

// Unsupported для неподдерживаемых методов
func (t *Transport) Unsupported(_ http.ResponseWriter, req *http.Request) {
	fmt.Printf("Unsupported method для server: %s", req.Method)
}
