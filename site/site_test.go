package site

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	testfileContent = []byte("https://www.cnn.com,covid\nhttps://www.google.com,\nhttps://aiven.io,kafka")
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
