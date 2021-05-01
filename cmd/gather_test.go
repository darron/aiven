package cmd

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	testfileContent = []byte("https://www.cnn.com,covid\nhttps://www.google.com,\nhttps://aiven.io,kafka")
)

func TestGetSites(t *testing.T) {
	f, err := testFile(t)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f)
	sites, err := getSites(f)
	if err != nil {
		t.Error(err)
	}
	want := 3
	if len(sites) != want {
		t.Errorf("getSites Want: %d, Got: %d", want, len(sites))
	}
}

func testFile(t *testing.T) (string, error) {
	tmp, err := ioutil.TempFile("", "aiven")
	if err != nil {
		t.Error(err)
	}
	_, err = tmp.Write(testfileContent)
	if err != nil {
		t.Error(err)
	}
	return tmp.Name(), tmp.Close()
}
