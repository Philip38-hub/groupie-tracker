package groupie

import (
	"encoding/json"
	"log"
	"html/template"
	"io"
	"net/http"
	"time"
)

// Define a struct to match the structure of the API response
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Members      []string `json:"members"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var error []string
	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 20 * time.Second, // 20-second timeout
	}

	// Make the GET request with the custom client
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Printf("Failed to get data from api: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
	defer resp.Body.Close()

	// Read and parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}
	// log.Printf(string(body))
	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Printf("Failed to parse JSON: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Load and parse the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Failed to open template index.html: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusInternalServerError, error)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, artists)
	if err != nil {
		log.Printf("Failed to execute template: %s", err)
		error = append(error, "Internal Server Error")
		ErrorHandler(w, r, http.StatusMethodNotAllowed, error)
		return
	}
}
