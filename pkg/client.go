package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

type Getter interface {
	Get(url string) (resp *http.Response, err error)
}

type EndOfLifeDateClient struct {
	Getter Getter
}

func NewEndOfLifeDateClient(getter Getter) *EndOfLifeDateClient {
	return &EndOfLifeDateClient{Getter: getter}
}

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
