// handlers/errors.go
package groupie

import (
	"html/template"
	"log"
	"net/http"
)

type ErrorData struct {
	Code   int
	Errors []string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, code int, errors []string) {
	w.WriteHeader(code)
	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		log.Printf("Failed to parse template: %s", err)
		return
	}

	data := ErrorData{Code: code, Errors: errors}
	tmpl.Execute(w, data)
}
