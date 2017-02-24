package macaque

import (
	"reflect"
	"testing"

	"github.com/mholt/caddy"
)

func TestConfigParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		expected  MacaqueConfig
	}{
		{
			`macaque / {
				policy keyless per_ip 100 5
				policy default per_ip 1000 5 per_key 50000 1440
				database /tmp/macaque.db
				sync_interval 10
				portal_api /portal_api
			}`, false, MacaqueConfig{
				Path: "/",
				Policies: map[string]MacaquePolicy{
					"keyless": MacaquePolicy{Name: "keyless", PerIPMax: 100, PerIPInterval: 5, PerKeyMax: 0, PerKeyInterval: 0},
					"default": MacaquePolicy{Name: "default", PerIPMax: 1000, PerIPInterval: 5, PerKeyMax: 50000, PerKeyInterval: 1440},
				},
				SyncInterval: 10,
				DBPath:       "/tmp/macaque.db",
				PortalAPI:    "/portal_api",
			},
		},
	}

	for i, test := range tests {
		actual, err := macaqueParseConfig(caddy.NewTestController("http", test.input))
		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("Test %d: Config was parsed incorrectly.\nExpected: %v\nGot     : %v", i, test.expected, actual)
		}
	}
}
