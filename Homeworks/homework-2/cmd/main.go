package main

import (
	server "homework/internal/pkg/server"
	serverWithData "homework/internal/pkg/server_with_data"
	"log"
	"net/http"
)

const (
	serverPort         = ":9000"
	serverWithDataPort = ":9001"
)

// Запускает 2 сервера: один простой и один, реализующий работу с данными запроса
func main() {
	go func() {
		runServerWithData()
	}()
	runServer()
}

// Простой сервер
func runServer() {
	implementation := server.MakeServer()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implementation.Create(res, req)
		case http.MethodGet:
			implementation.Read(res, req)
		case http.MethodPut:
			implementation.Update(res, req)
		case http.MethodDelete:
			implementation.Delete(res, req)
		default:
			implementation.Unsupported(res, req)
		}
	})
	if err := http.ListenAndServe(serverPort, mux); err != nil {
		log.Fatal(err)
	}
}

// Сервер, реализующий базовую работу с данными запроса
func runServerWithData() {
	implementation := serverWithData.MakeServer()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			implementation.Create(res, req)
		case http.MethodGet:
			implementation.Read(res, req)
		case http.MethodPut:
			implementation.Update(res, req)
		case http.MethodDelete:
			implementation.Delete(res, req)
		default:
			implementation.Unsupported(res, req)
		}
	})
	if err := http.ListenAndServe(serverWithDataPort, mux); err != nil {
		log.Fatal(err)
	}
}
