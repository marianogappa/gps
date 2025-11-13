package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const cacheFileName = "city_cache.json"

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
