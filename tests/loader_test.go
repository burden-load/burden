package tests

import (
	"burden/internal/loader"
	"testing"
)

func TestLoadCollection(t *testing.T) {
	requests, err := loader.LoadCollection("../examples/example_collection.json")
	if err != nil {
		t.Fatalf("Error loading collection: %v", err)
	}

	if len(requests) != 21 {
		t.Errorf("Expected 2 requests, got %d", len(requests))
	}
}
