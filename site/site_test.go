package site

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	testfileContent = []byte("https://www.cnn.com,covid\nhttps://www.google.com,\nhttps://aiven.io,kafka")
	testMetricJSON  = []byte(`{"captured_at":"2021-05-03T17:59:49.124364972Z","address":"https://aiven.io","response_time":284510084,"status_code":200,"regexp":"kafka","regexp_status":true}`)
)

func TestGetEntries(t *testing.T) {
	f, err := testFile(t, testfileContent)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f)
	sites, err := GetEntries(f)
	if err != nil {
		t.Error(err)
	}
	want := 3
	if len(sites) != want {
		t.Errorf("GetEntries Want: %d, Got: %d", want, len(sites))
	}
}

func TestGetEntriesFail(t *testing.T) {
	// This should fail - as it doesn't exist.
	_, err := GetEntries("filename-does-not-exist")
	if err == nil {
		t.Error("That should have failed")
	}
	// Setup a blank file.
	f, err := testFile(t, []byte(""))
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f)
	// This should fail as there are no entries.
	_, err = GetEntries(f)
	if err == nil {
		t.Error("That should have failed")
	}
}

func TestExtractMetrics(t *testing.T) {
	m, err := ExtractMetrics(testMetricJSON)
	if err != nil {
		t.Error(err)
	}
	if m.Address != "https://aiven.io" {
		t.Errorf("Want: %s Got: %s", "https://aiven.io", m.Address)
	}
	_, err = ExtractMetrics([]byte(""))
	if err == nil {
		t.Error("That should have errored")
	}
}

func TestGetMetrics(t *testing.T) {
	g := Entry{
		Address: "https://www.google.com/",
		Regexp:  "fake",
	}
	// This should work properly
	m, err := g.GetMetrics(1*time.Second, MockGoogleWorkingEntry{}, false)
	if err != nil {
		t.Error(err)
	}
	if m.StatusCode != 200 {
		t.Errorf("Want: %d, Got: %d", 200, m.StatusCode)
	}
	// This request should be broken
	m, err = g.GetMetrics(1*time.Second, MockGoogleBrokenEntry{}, false)
	if err == nil {
		t.Error("This should be broken")
	}
	if m.StatusCode != 0 {
		t.Error("That should have been blank")
	}
	// This request should fail because it exceeded the context timeout.
	_, err = g.GetMetrics(10*time.Millisecond, MockGoogleSlowEntry{}, false)
	if err == nil {
		t.Error("That should have failed")
	}
	if err.Error() != "context deadline exceeded" {
		t.Error("That should be a context error")
	}
}

func testFile(t *testing.T, content []byte) (string, error) {
	tmp, err := ioutil.TempFile("", "aiven")
	if err != nil {
		t.Error(err)
	}
	_, err = tmp.Write(content)
	if err != nil {
		t.Error(err)
	}
	return tmp.Name(), tmp.Close()
}
