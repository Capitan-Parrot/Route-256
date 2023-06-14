package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	headerKey = "hw-sum"
)

type server struct{}

// Create выводит хэдеры POST запроса и сумму headerKey + 5
func (s *server) Create(_ http.ResponseWriter, req *http.Request) {
	fmt.Printf("Create, headers: [%v]\n", req.Header)
	valueFromRequest := req.Header.Get(headerKey)
	if valueFromRequest == "" {
		fmt.Println("No hw-sum header")
		return
	}
	convertedValueFromRequest, err := strconv.Atoi(valueFromRequest)
	if err != nil {
		fmt.Printf("Incorrect hw-sum header: %s", err.Error())
		return
	}
	fmt.Printf("Sum: [%d]\n", convertedValueFromRequest+5)
}

// Read выводит параметры GET запроса
func (s *server) Read(_ http.ResponseWriter, req *http.Request) {
	fmt.Printf("Get, query params: [%v]\n", req.URL.Query())
}

// Update выводит тело PUT запроса
func (s *server) Update(_ http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("error while updating, err: %v\n", err.Error())
		return
	}
	fmt.Printf("Update, body: [%s]\n", string(body))
}

// Delete для DELETE метода
func (s *server) Delete(_ http.ResponseWriter, _ *http.Request) {
	fmt.Println("Delete")
}

// Unsupported для неподдерживаемых методов
func (s *server) Unsupported(_ http.ResponseWriter, req *http.Request) {
	fmt.Printf("Unsupported method для server: %s", req.Method)
}

func MakeServer() server {
	return server{}
}
