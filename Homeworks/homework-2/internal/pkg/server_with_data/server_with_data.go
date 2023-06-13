package server_with_data

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	idKey = "id"
)

type serverWithData struct {
	data map[uint32]string
}

type request struct {
	Id    uint32 `json:"id"`
	Value string `json:"value"`
}

// Create - создание записи на сервере в формате ключ (uint32) -значение (string)
func (s *serverWithData) Create(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error while reading request body: [%s]\n", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var unmarshalledRequest request
	if err = json.Unmarshal(body, &unmarshalledRequest); err != nil {
		log.Printf("Error while unmarshalling request body: [%s]\n", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, ok := s.data[unmarshalledRequest.Id]; ok {
		log.Println("Id is already existed")
		res.WriteHeader(http.StatusConflict)
		return
	}
	if unmarshalledRequest.Value == "" {
		log.Println("No value key")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	s.data[unmarshalledRequest.Id] = unmarshalledRequest.Value
}

// Read - чтение данных (string) с сервера по ключу (uint32)
func (s *serverWithData) Read(res http.ResponseWriter, req *http.Request) {
	idFromRequest := req.URL.Query().Get(idKey)
	if idFromRequest == "" {
		log.Println("No id key")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	convertedIdFromRequest, err := strconv.ParseUint(idFromRequest, 10, 32)
	if err != nil {
		log.Printf("Incorrect id key: %s", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	valueFromData, ok := s.data[uint32(convertedIdFromRequest)]
	if !ok {
		log.Printf("No data for id: %s", idFromRequest)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err := res.Write([]byte(valueFromData)); err != nil {
		log.Printf("Read method of serverWithData: %s", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Update - Обновление данных (string) по существующему ключу (uint32)
func (s *serverWithData) Update(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error while reading request body: [%s]", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unmarshalledRequest request
	if err = json.Unmarshal(body, &unmarshalledRequest); err != nil {
		log.Printf("Error while unmarshalling request body: [%s]", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, ok := s.data[unmarshalledRequest.Id]; !ok {
		log.Printf("No data for id: %d", unmarshalledRequest.Id)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if unmarshalledRequest.Value == "" {
		log.Printf("No value key")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	s.data[unmarshalledRequest.Id] = unmarshalledRequest.Value
}

// Delete - удаление данных по существующему ключу (uint32)
func (s *serverWithData) Delete(res http.ResponseWriter, req *http.Request) {
	idFromRequest := req.URL.Query().Get(idKey)
	if idFromRequest == "" {
		log.Println("No id key")
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	convertedIdFromRequest, err := strconv.ParseUint(idFromRequest, 10, 32)
	if err != nil {
		log.Printf("Incorrect id key: %s", err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, ok := s.data[uint32(convertedIdFromRequest)]
	if !ok {
		log.Printf("No data for id: %s", idFromRequest)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	delete(s.data, uint32(convertedIdFromRequest))
}

// Unsupported для неподдерживаемых методов
func (s *serverWithData) Unsupported(_ http.ResponseWriter, req *http.Request) {
	fmt.Printf("Unsupported method for serverWithData: %s", req.Method)
}

func MakeServer() serverWithData {
	return serverWithData{data: make(map[uint32]string)}
}
