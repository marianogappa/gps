package pkg

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

const (
	nominatimAPIURL = "https://nominatim.openstreetmap.org/search?format=json&limit=1&q="
)

func fetchLocation(searchTerm string, httpClient http.Client) (Location, error) {
	// Encode search term
	encodedSearchTerm := url.QueryEscape(searchTerm)

	// Make API request
	resp, err := httpClient.Get(nominatimAPIURL + encodedSearchTerm)
	if err != nil {
		return Location{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Location{}, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Decode API response
	var locations []Location
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return Location{}, err
	}

	if len(locations) == 0 {
		return Location{}, fmt.Errorf("no results found for %s", searchTerm)
	}

	return locations[0], nil
}

func getHttpClient(gpsOpts GPSOptions) http.Client {
	var dialer proxy.Dialer = &net.Dialer{}
	if gpsOpts.Dialer != nil {
		dialer = gpsOpts.Dialer
	}
	return http.Client{Transport: &http.Transport{
		Dial: dialer.Dial,
	}}
}
