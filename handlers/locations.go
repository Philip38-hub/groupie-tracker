package groupie

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Define a struct to match the structure of the API response
type DateLocation struct {
	ID int `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func LocationsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Make the GET request to fetch location data
	resp, err := client.Get("https://groupietrackers.herokuapp.com/api/locations") // Update with correct URL
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

	var locations map[string]DateLocation
	err = json.Unmarshal(body, &locations)
	if err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	// Find the location data for the requested artist ID
	locationData, ok := locations[artistID]
	if !ok {
		http.Error(w, "Artist ID not found", http.StatusNotFound)
		return
	}

	// Return the location data as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locationData); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
