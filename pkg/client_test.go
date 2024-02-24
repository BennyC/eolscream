package pkg

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockGetter struct {
	Response *http.Response
	Error    error
}

func (m *MockGetter) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

func TestEndOfLifeDateClient_FetchProductInfo_Success(t *testing.T) {
	client := NewEndOfLifeDateClient(&MockGetter{})

	tests := []struct {
		name            string
		mockGetter      *MockGetter
		product         Product
		wantReleaseInfo *ReleaseInfo
	}{
		{
			name: "ValidProduct",
			mockGetter: &MockGetter{
				Response: &http.Response{
					Body: io.NopCloser(strings.NewReader(`{
						"releaseDate": "2023-05-23",
						"eol": "2026-04-01"
					}`)),
				},
			},
			product: Product{Name: "mockProduct", Version: "0.0.1"},
			wantReleaseInfo: &ReleaseInfo{
				ReleaseDate:   "2023-05-23",
				EndOfLifeDate: "2026-04-01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.Getter = tt.mockGetter
			got, err := client.FetchProductInfo(tt.product)

			assert.Equal(t, tt.wantReleaseInfo, got)
			assert.Nil(t, err)
		})
	}
}

func TestEndOfLifeDateClient_FetchProductInfo_Errors(t *testing.T) {
	client := NewEndOfLifeDateClient(&MockGetter{})

	tests := []struct {
		name       string
		mockGetter *MockGetter
		product    Product
		wantErr    error
	}{
		{
			name: "ErrorGettingProductInfo",
			mockGetter: &MockGetter{
				Error: errors.New("mock error"),
			},
			product: Product{Name: "mockProduct", Version: "0.0.1"},
			wantErr: errors.New("mock error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.Getter = tt.mockGetter
			got, err := client.FetchProductInfo(tt.product)

			assert.Nil(t, got)
			assert.NotNil(t, err)
			assert.Equal(t, tt.wantErr.Error(), err.Error())
		})
	}
}
