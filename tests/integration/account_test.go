package integration_test

import (
	"testing"
)

func TestAccountDeserialization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping account integration test in short mode.")
	}
	// TODO: Add integration test for fetching and deserializing raw accounts
}
