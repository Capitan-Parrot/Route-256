package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"workshop.4.2/internal/model"
	usercore "workshop.4.2/internal/pkg/service/user"
)

const (
	addr = ":9000"
)

func main() {
	server := transport{
		server: usercore.New(),
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
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(addr, mux))
}

type transport struct {
	server *usercore.Server
}

func (t *transport) Get(w http.ResponseWriter, r *http.Request) {
	value := r.URL.Query().Get(model.QueryParamUserID)
	userID, err := strconv.ParseUint(value, 10, 64)
	if err != nil || userID == 0 {
		log.Printf("invalid user.id, url: %s, err: %v", r.URL, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := t.server.Get(model.ClientID(userID))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(result); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
