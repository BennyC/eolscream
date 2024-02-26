package pkg

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

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
