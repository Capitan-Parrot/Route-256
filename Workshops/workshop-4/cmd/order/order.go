package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"workshop.4.2/internal/model"
	ordercore "workshop.4.2/internal/pkg/service/order"
)

const (
	addr = ":9001"
)

func main() {
	server := transport{
		server: ordercore.New(),
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch r.Method {
		case http.MethodGet:
			server.Get(w, r)
		case http.MethodPost:
			server.Create(w, r)
		case http.MethodPut:
			server.Update(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(addr, mux))
}

type transport struct {
	server *ordercore.Server
}

func (t *transport) Get(w http.ResponseWriter, r *http.Request) {
	queryOrderID := r.URL.Query().Get(model.QueryParamOrderID)
	orderID, err := strconv.ParseUint(queryOrderID, 10, 64)
	if err != nil || orderID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := t.server.Get(model.OrderID(orderID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *transport) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var order model.Order
	if err = json.Unmarshal(body, &order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = t.server.Update(order); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (t *transport) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var order model.Order
	if err = json.Unmarshal(body, &order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = t.server.Create(order); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
