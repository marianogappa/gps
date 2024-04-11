package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	nominatimAPIURL = "https://nominatim.openstreetmap.org/search?format=json&limit=1&q="
	cacheFileName   = "city_cache.json"
)

type Location struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func main() {
	// Parse command-line flags
	separator := flag.String("separator", "tab", "Separator for output lines: tab or comma")
	flag.Parse()

	// Load cache file
	cache := loadCache()

	// Read search terms from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		searchTerm := scanner.Text()

		// Check cache for existing data
		if location, ok := cache[searchTerm]; ok {
			printLocation(searchTerm, location, *separator)
			continue
		}

		// Fetch data from API
		location, err := fetchLocation(searchTerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching data for %s: %v\n", searchTerm, err)
			continue
		}

		// Update cache
		cache[searchTerm] = location
		saveCache(cache)

		// Print location
		printLocation(searchTerm, location, *separator)
	}
}

func fetchLocation(searchTerm string) (Location, error) {
	// Encode search term
	encodedSearchTerm := url.QueryEscape(searchTerm)

	// Make API request
	resp, err := http.Get(nominatimAPIURL + encodedSearchTerm)
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
		return Location{}, fmt.Errorf("No results found for %s", searchTerm)
	}

	return locations[0], nil
}

func printLocation(searchTerm string, location Location, separator string) {
	displayName := strings.ReplaceAll(searchTerm, ",", "") // Use user's chosen search term
	switch separator {
	case "comma":
		fmt.Printf("%s,%s,%s\n", location.Lat, location.Lon, displayName)
	default:
		fmt.Printf("%s\t%s\t%s\n", location.Lat, location.Lon, displayName)
	}
}

func loadCache() map[string]Location {
	cache := make(map[string]Location)

	// Open cache file
	cacheFilePath := filepath.Join(os.TempDir(), cacheFileName)
	file, err := os.Open(cacheFilePath)
	if err != nil {
		return cache
	}
	defer file.Close()

	// Read cache data
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return cache
	}

	// Unmarshal JSON data
	if err := json.Unmarshal(data, &cache); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling cache data: %v\n", err)
	}

	return cache
}

func saveCache(cache map[string]Location) {
	// Marshal cache data to JSON
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling cache data: %v\n", err)
		return
	}

	// Write data to cache file
	cacheFilePath := filepath.Join(os.TempDir(), cacheFileName)
	if err := ioutil.WriteFile(cacheFilePath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing cache data to file: %v\n", err)
	}
}
