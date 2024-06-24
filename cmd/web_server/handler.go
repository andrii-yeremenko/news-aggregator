package web_server

import "net/http"

// Handler is a function that is responsible for managing or dealing with specific events, actions, or errors that occur
// during the execution of web server.
type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
