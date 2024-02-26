package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct {
	expectedReleaseInfo *ReleaseInfo
	expectedErr         error
}

func (mc *MockClient) FetchReleaseInfo(product Product) (*ReleaseInfo, error) {
	return mc.expectedReleaseInfo, mc.expectedErr
}

type MockNotifier struct {
	notifyCount int
}

func (m *MockNotifier) Notify(p Product, r ReleaseInfo) {
	m.notifyCount++
}

func TestNotifyNearEndOfLife(t *testing.T) {
	testCases := []struct {
		name                string
		path                string
		expectedReleaseInfo *ReleaseInfo
		expectedNotifyCount int
	}{
		{
			name:                "happy",
			path:                "testdata/good_catalogue.json",
			expectedReleaseInfo: &ReleaseInfo{ReleaseDate: "2019-05-01", EndOfLifeDate: "2024-02-26"},
			expectedNotifyCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			notifier := &MockNotifier{notifyCount: 0}
			opts := CatalogueCheckerOptions{
				Path: tc.path,
				Client: &MockClient{
					expectedReleaseInfo: tc.expectedReleaseInfo,
				},
				Notifier: notifier,
			}

			cc := NewCatalogueChecker(opts)
			err := cc.NotifyNearEndOfLife()

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedNotifyCount, notifier.notifyCount)
		})
	}
}
