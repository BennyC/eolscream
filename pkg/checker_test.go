package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	expectedReleaseInfo *ReleaseInfo
	expectedErr         error
}

func (mc *mockClient) FetchReleaseInfo(product Product) (*ReleaseInfo, error) {
	return mc.expectedReleaseInfo, mc.expectedErr
}

type mockNotifier struct {
	notifyCount int
}

func (m *mockNotifier) Notify(p Product, r ReleaseInfo) {
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

		{
			name:                "happy in future",
			path:                "testdata/good_catalogue.json",
			expectedReleaseInfo: &ReleaseInfo{ReleaseDate: "2019-05-01", EndOfLifeDate: "2100-02-26"},
			expectedNotifyCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			notifier := &mockNotifier{notifyCount: 0}
			opts := CatalogueCheckerOptions{
				Path: tc.path,
				Client: &mockClient{
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
