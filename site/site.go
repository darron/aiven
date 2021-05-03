package site

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gojektech/heimdall/v6"
)

// Metrics is the data we pass along via Kafka to be stored in Postgres
type Metrics struct {
	CapturedAt   time.Time     `json:"captured_at"`
	Address      string        `json:"address"`
	ResponseTime time.Duration `json:"response_time"`
	StatusCode   int           `json:"status_code"`
	Regexp       string        `json:"regexp"`
	RegexpStatus bool          `json:"regexp_status"`
}

// Entries is a slice of Entry values
type Entries []Entry

// Entry represents an Address and optional Regexp to query for HTTP metrics.
type Entry struct {
	Address string
	Regexp  string
}

// ExtractMetrics takes a byte string of JSON and converts to Metrics
func ExtractMetrics(j []byte) (Metrics, error) {
	var m Metrics
	err := json.Unmarshal(j, &m)
	return m, err
}

// GetMetrics returns metrics data for an Entry.
// TODO: Setup mocks for better tests.
func (e Entry) GetMetrics(timeout time.Duration, h heimdall.Doer, debug bool) (Metrics, error) {
	var m Metrics

	// Let's set the stuff we know already.
	m.Address = e.Address
	m.Regexp = e.Regexp
	m.CapturedAt = time.Now().UTC()

	// Set a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Setup the request.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, e.Address, nil)
	if err != nil {
		return m, err
	}

	if debug {
		fmt.Printf("%#v\n", req)
	}

	// Let's do the request.
	start := time.Now()
	res, err := h.Do(req)
	if err != nil {
		return m, err
	}

	// We don't need the body unless we've got a regexp.
	if e.Regexp != "" {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return m, err
		}
		// Do the regexp here and assign the value.
		match, _ := regexp.MatchString(e.Regexp, string(body))
		m.RegexpStatus = match
	}
	took := time.Since(start)

	// Set the last few values
	m.StatusCode = res.StatusCode
	m.ResponseTime = took

	if debug {
		fmt.Printf("Result: %#v\n", res)
		fmt.Printf("Time Taken: %s\n", took)
		fmt.Printf("Metrics: %#v\n", m)
	}

	if ctx.Err() != nil {
		return m, errors.New("context deadline exceeded")
	}
	return m, nil
}

// GetEntries reads a file and returns website Entries.
func GetEntries(filename string) (Entries, error) {
	var entries Entries
	f, err := os.Open(filename)
	if err != nil {
		return entries, err
	}
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		return entries, err
	}
	for _, line := range lines {
		s := Entry{
			Address: line[0],
			// TODO: Add regexp checking to make sure that it's a regexp - maybe.
			Regexp: line[1],
		}
		// Only add this if we have an address - skip otherwise.
		if s.Address != "" {
			entries = append(entries, s)
		}
	}
	if len(entries) == 0 {
		return entries, errors.New("no entries - cannot proceed")
	}
	return entries, nil
}
