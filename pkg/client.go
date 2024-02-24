package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ReleaseInfo struct {
	ReleaseDate   string `json:"releaseDate"`
	EndOfLifeDate string `json:"eol"`
}

type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

type EndOfLifeDateClient struct {
	Getter Getter
}

func NewEndOfLifeDateClient(getter Getter) *EndOfLifeDateClient {
	return &EndOfLifeDateClient{Getter: getter}
}

// FetchProductInfo is a method of the EndOfLifeDateClient struct that fetches product information using an HTTP GET request.
// It takes a Product as input and returns the ReleaseInfo and an error (if any).
func (c *EndOfLifeDateClient) FetchProductInfo(product Product) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://endoflife.date/api/%s/%s", product.Name, product.Version)
	resp, err := c.Getter.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return extractReleaseInfo(resp)
}

func extractReleaseInfo(resp *http.Response) (*ReleaseInfo, error) {
	var releaseInfo ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&releaseInfo); err != nil {
		return nil, err
	}
	return &releaseInfo, nil
}
