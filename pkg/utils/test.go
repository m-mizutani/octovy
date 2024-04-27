package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// LoadEnv loads environment variable and return its value for test code. If the variable is not set, it skips the test.
func LoadEnv(t *testing.T, name string) string {
	t.Helper()
	v, ok := os.LookupEnv(name)
	if !ok {
		t.Skipf("Skip test because %s is not set", name)
	}

	return v
}

// LoadJson loads JSON file and decode it into v for test code. If it fails, it fails the test.
func LoadJson(t *testing.T, path string, v interface{}) {
	t.Helper()
	fp, err := os.Open(filepath.Clean(path))
	if err != nil {
		t.Fatalf("Failed to open file: %s", path)
	}
	defer SafeClose(fp)

	if err := json.NewDecoder(fp).Decode(v); err != nil {
		t.Fatalf("Failed to decode JSON: %s", err)
	}
}

// DumpJson dumps v into JSON file for test code. If it fails, it fails the test.
func DumpJson(t *testing.T, path string, v interface{}) {
	if t != nil {
		t.Helper()
	}
	fp, err := os.Create(filepath.Clean(path))
	if err != nil {
		if t != nil {
			t.Fatalf("Failed to create file: %s", path)
		}
		panic(err)
	}
	defer SafeClose(fp)

	if err := json.NewEncoder(fp).Encode(v); err != nil {
		if t != nil {
			t.Fatalf("Failed to encode JSON: %s", err)
		}
		panic(err)
	}
}
