package groupie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Struct to hold the dates data
type Dates struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

func DatesHandler(w http.ResponseWriter, r *http.Request) {
	// Get the artist ID from the query parameters
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		http.Error(w, "Missing artist ID", http.StatusBadRequest)
		return
	}

	// Create a custom HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second, // 10-second timeout
	}

	// Make the GET request to fetch dates data
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/dates") // Update with the correct URL
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

	var dates Dates
	err = json.Unmarshal(body, &dates)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Find the dates data for the requested artist ID
	var datesData struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}
	found := false
	for _, date := range dates.Index {
		id, err := strconv.Atoi(artistID)
		if err != nil {
			http.Error(w, "Invalid artist ID", http.StatusBadRequest)
			return
		}
		if date.ID == id {
			datesData = date
			found = true
			break
		}
	}

	// If the artist ID is not found, return an error
	if !found {
		http.Error(w, "Artist ID not found", http.StatusNotFound)
		return
	}
	fmt.Println(datesData)
	// Return the dates data as JSON
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8080")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(datesData); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
