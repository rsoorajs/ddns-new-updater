package dns

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newDefaultSettings(t *testing.T) {
	t.Parallel()

	settings := newDefaultSettings()

	assert.NotEmpty(t, settings.providers)
	assert.GreaterOrEqual(t, int(settings.timeout), int(time.Millisecond))
}

func Test_SetProviders(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  settings
		providers        []Provider
		expectedSettings settings
		err              error
	}{
		"OpenDNS": {
			initialSettings: settings{
				providers: []Provider{Cloudflare},
			},
			providers: []Provider{OpenDNS},
			expectedSettings: settings{
				providers: []Provider{OpenDNS},
			},
		},
		"OpenDNS and Cloudflare": {
			initialSettings: settings{
				providers: []Provider{Cloudflare},
			},
			providers: []Provider{OpenDNS, Cloudflare},
			expectedSettings: settings{
				providers: []Provider{Cloudflare, OpenDNS},
			},
		},
		"invalid provider": {
			initialSettings: settings{
				providers: []Provider{Cloudflare},
			},
			providers: []Provider{Provider("invalid")},
			expectedSettings: settings{
				providers: []Provider{Cloudflare},
			},
			err: errors.New("unknown public IP echo DNS provider: invalid"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := testCase.initialSettings

			option := SetProviders(testCase.providers[0], testCase.providers[1:]...)
			err := option(&settings)

			assert.Equal(t, testCase.expectedSettings, settings)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_SetTimeout(t *testing.T) {
	t.Parallel()

	initialSettings := settings{}
	expectedSettings := settings{
		timeout: time.Hour,
	}

	option := SetTimeout(time.Hour)
	err := option(&initialSettings)

	require.NoError(t, err)
	assert.Equal(t, expectedSettings, initialSettings)
}
