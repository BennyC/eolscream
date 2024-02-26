package pkg

import (
	"os"
	"testing"
	"time"

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

func TestReleaseInfo_IsEndOfLifeDateNear(t *testing.T) {
	tests := []struct {
		name          string
		endOfLifeDate string
		want          bool
		wantErr       bool
	}{
		{
			name:          "Given date is within six months from today",
			endOfLifeDate: time.Now().AddDate(0, 5, 0).Format("2006-01-02"),
			want:          true,
			wantErr:       false,
		},
		{
			name:          "Given date is more than six months from today",
			endOfLifeDate: time.Now().AddDate(0, 7, 0).Format("2006-01-02"),
			want:          false,
			wantErr:       false,
		},
		{
			name:          "Given date is in the past",
			endOfLifeDate: time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
			want:          true,
			wantErr:       false,
		},
		{
			name:          "Given date is in a wrong format",
			endOfLifeDate: "202-02-900",
			want:          false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReleaseInfo{EndOfLifeDate: tt.endOfLifeDate}
			got, err := r.IsEndOfLifeDateNear()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
