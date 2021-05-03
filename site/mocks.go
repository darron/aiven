package site

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type MockGoogleWorkingEntry struct{}

func (m MockGoogleWorkingEntry) Do(*http.Request) (*http.Response, error) {
	body := `This is the very fake body of a Google request`
	res := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
	return &res, nil
}

type MockGoogleBrokenEntry struct{}

func (m MockGoogleBrokenEntry) Do(*http.Request) (*http.Response, error) {
	body := `This should not work`
	res := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
	return &res, errors.New("this is b0rken")
}

type MockGoogleSlowEntry struct{}

func (m MockGoogleSlowEntry) Do(*http.Request) (*http.Response, error) {
	body := `This should not work - context tieout`
	res := http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
	// This should put it past the context timeout.
	time.Sleep(100 * time.Millisecond)
	return &res, nil
}
