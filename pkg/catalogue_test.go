package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatalogueFromJson(t *testing.T) {
	// Test data
	data := []byte(`{"catalogue": [{"name": "A", "version": "1"}, {"name": "B", "version": "2"}]}`)

	catalogue, err := CatalogueFromJson(data)

	assert.NoError(t, err, "Expected no error for valid json")
	assert.Len(t, *catalogue, 2, "Expected 2 items in catalogue")
}

func TestNewCatalogueFromFile(t *testing.T) {
	testFile := "testFile.json"
	testData := []byte(`{"catalogue": [{"name": "A", "version": "1"}, {"name": "B", "version": "2"}]}`)

	// Write test data to file
	err := os.WriteFile(testFile, testData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer func(name string) {
		_ = os.Remove(name)
	}(testFile)

	// Call the function we want to test
	catalogue, err := NewCatalogueFromFile(testFile)

	assert.NoError(t, err, "Expected no error for valid file")
	assert.Len(t, *catalogue, 2, "Expected 2 items in catalogue")
}
