package pkg

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/proxy"
)

type GPSOptions struct {
	OutputSeparator string
	OutputToWriter  io.Writer
	OutputToChannel chan string
	Dialer          proxy.Dialer
}

func WithOutputSeparator(separator string) func(*GPSOptions) {
	return func(opts *GPSOptions) {
		opts.OutputSeparator = separator
	}
}

func WithOutputToWriter(w io.Writer) func(*GPSOptions) {
	return func(opts *GPSOptions) {
		opts.OutputToWriter = w
	}
}

func WithOutputChannel(ch chan string) func(*GPSOptions) {
	return func(opts *GPSOptions) {
		opts.OutputToChannel = ch
	}
}

func WithDialer(d proxy.Dialer) func(*GPSOptions) {
	return func(opts *GPSOptions) {
		opts.Dialer = d
	}
}

func GPSFromReader(r io.Reader, opts ...func(*GPSOptions)) {
	gpsOpts := GPSOptions{OutputSeparator: "tab"}
	for _, opt := range opts {
		opt(&gpsOpts)
	}

	httpClient := getHttpClient(gpsOpts)

	// Load cache file
	cache := loadCache()

	// Read search terms from reader
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		searchTerm := scanner.Text()

		// Check cache for existing data
		if location, ok := cache[searchTerm]; ok {
			outputLocation(searchTerm, location, gpsOpts)
			continue
		}

		// Fetch data from API
		location, err := fetchLocation(searchTerm, httpClient)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching data for %s: %v\n", searchTerm, err)
			continue
		}

		// Update cache
		cache[searchTerm] = location
		saveCache(cache)

		// Print location
		outputLocation(searchTerm, location, gpsOpts)
	}
}

func prepareOutput(searchTerm string, location Location, separator string) string {
	displayName := strings.ReplaceAll(searchTerm, ",", "") // Use user's chosen search term
	switch separator {
	case "comma":
		return fmt.Sprintf("%s,%s,%s", location.Lat, location.Lon, displayName)
	default:
		return fmt.Sprintf("%s\t%s\t%s", location.Lat, location.Lon, displayName)
	}
}

func outputLocation(searchTerm string, location Location, gpsOpts GPSOptions) {
	s := prepareOutput(searchTerm, location, gpsOpts.OutputSeparator)
	if gpsOpts.OutputToChannel != nil {
		gpsOpts.OutputToChannel <- s
	}
	if gpsOpts.OutputToWriter != nil {
		fmt.Fprintln(gpsOpts.OutputToWriter, s)
	}
}
