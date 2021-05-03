package cmd

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Clear the environment - Load should fail.
	os.Clearenv()
	_, err := Load("gather")
	if err == nil {
		t.Error("That should have failed")
	}
	// Set the only thing we need for "gather"
	os.Setenv("KAFKA_HOST", "not.a.real.domain.name:12345")
	_, err = Load("gather")
	if err != nil {
		t.Error("That should NOT have failed")
	}
	// Set what we need for "store"
	os.Setenv("POSTGRES_URL", "postgres://user:password@not.real.domain.name:20349/defaultdb?sslmode=disable")
	_, err = Load("store")
	if err != nil {
		t.Error("That should NOT have failed")
	}
	// This should fail now.
	_, err = Load("not-a-known-configType")
	if err.Error() != "must specify known configType" {
		t.Error("That should have errored")
	}
}
