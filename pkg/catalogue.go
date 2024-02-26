package pkg

import (
	"encoding/json"
	"os"
	"time"
)

type ProductFile struct {
	Catalogue `json:"catalogue"`
}

type Product struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Catalogue []Product

// NewCatalogueFromFile reads the contents of a file at the specified path,
// then calls CatalogueFromJson to parse the contents into a Catalogue object.
// It returns a pointer to the resulting Catalogue and any error encountered.
// Example usage:
//
//	catalogue, err := NewCatalogueFromFile("/path/to/file.json")
func NewCatalogueFromFile(path string) (*Catalogue, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return CatalogueFromJson(contents)
}

// CatalogueFromJson takes a byte slice representing JSON data and returns a pointer to a Catalogue object and any error encountered during parsing.
// It uses the ProductFile struct to unmarshal the JSON data and retrieves the Catalogue field.
// Example usage:
//
//	bytes := []byte(`{"catalogue": [{"name": "A", "version": "1"}, {"name": "B", "version": "2"}]}`)
//	catalogue, err := CatalogueFromJson(bytes)
//	if err != nil {
//	  fmt.Println("Error:", err)
//	  return
//	}
//	fmt.Println("Catalogue:", *catalogue)
//
// Note: The Catalogue object is defined as []Product, where Product is a struct with Name and Version fields.
//
//	The ProductFile struct is defined with a field named Catalogue that will be used for unmarshaling.
//	The CatalogueFromJson function can be used in conjunction with NewCatalogueFromFile to read JSON data from a file and obtain a Catalogue object.
//	There is also a test case TestCatalogueFromJson that demonstrates the usage of CatalogueFromJson with sample data.
func CatalogueFromJson(bytes []byte) (*Catalogue, error) {
	var products ProductFile
	err := json.Unmarshal(bytes, &products)
	if err != nil {
		return nil, err
	}

	return &products.Catalogue, nil
}

type ReleaseInfo struct {
	ReleaseDate   string `json:"releaseDate"`
	EndOfLifeDate string `json:"eol"`
}

func (r *ReleaseInfo) IsEndOfLifeDateNear() (bool, error) {
	eol, err := time.Parse("2006-01-02", r.EndOfLifeDate)
	if err != nil {
		return false, err
	}

	today := time.Now()
	sixMonthsFromNow := today.AddDate(0, 6, 0)

	return eol.Before(sixMonthsFromNow), nil
}
