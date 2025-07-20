package http

import (
	"net/http"

	"github.com/go-chi/render"
)

// SendMessage is a utility function to send a JSON response with a message.
// It sets the status code and writes a JSON response with the provided message.
func SendMessage(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	render.Status(r, statusCode)
	render.JSON(w, r, map[string]string{"message": message})
}
