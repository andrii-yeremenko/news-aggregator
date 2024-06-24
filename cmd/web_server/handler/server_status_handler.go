package handler

import (
	"fmt"
	"net/http"
)

type ServerStatusHandler struct {
}

func NewServerStatusHandler() *ServerStatusHandler {
	return &ServerStatusHandler{}
}

func (ssh *ServerStatusHandler) Handle(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "It Works!")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
