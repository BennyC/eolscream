package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

type EndOfLifeClient interface {
	FetchReleaseInfo(product Product) (*ReleaseInfo, error)
}

type EndOfLifeHttpClient struct {
	Getter Getter
}

func NewEndOfLifeDateClient(getter Getter) *EndOfLifeHttpClient {
	return &EndOfLifeHttpClient{Getter: getter}
}

func (c *EndOfLifeHttpClient) FetchReleaseInfo(product Product) (*ReleaseInfo, error) {
	url := fmt.Sprintf("https://endoflife.date/api/%s/%s.json", product.Name, product.Version)
	resp, err := c.Getter.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return extractReleaseInfo(resp)
}

func extractReleaseInfo(resp *http.Response) (*ReleaseInfo, error) {
	var releaseInfo ReleaseInfo
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &releaseInfo); err != nil {
		return nil, err
	}
	return &releaseInfo, nil
}
