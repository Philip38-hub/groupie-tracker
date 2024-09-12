package groupie

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"time"
)

// Define a struct to match the structure of the API response
type Artist struct {
	ID 			int `json:"id"`
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
	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // 10-second timeout
	}

	// Make the GET request with the custom client
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read and parse the JSON response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	// fmt.Println(string(body))
	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Load and parse the template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the data
	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
