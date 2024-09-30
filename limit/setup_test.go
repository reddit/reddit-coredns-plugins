package limit

import (
	"github.com/coredns/caddy"
	"strings"
	"testing"
)

// Errors that come from our plugin should start with this string
const errPrefix = "plugin/limit:"

func TestSetuplimit(t *testing.T) {
	tests := []struct {
		input           string
		shouldErr       bool
		expectedRecords int
	}{
		{`limit`, true, -1},
		{`limit "0"`, true, -1},
		{`limit "1"`, false, 1},
		{`limit "9000"`, false, 9000},
		{`limit "512 512"`, true, -1},
		{`limit "abc123"`, true, -1},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		limit, err := parse(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error but found %s for input %s", i, err, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: Error found for input %s. Error: %v", i, test.input, err)
			}

			if !strings.Contains(err.Error(), errPrefix) {
				t.Errorf("Test %d: Expected error to contain: %v, found error: %v, input: %s", i, errPrefix, err, test.input)
			}
		}

		if !test.shouldErr && limit != test.expectedRecords {
			t.Errorf("Test %d: limit not correctly set for input %s. Expected: %d, actual: %d", i, test.input, test.expectedRecords, limit)
		}
	}
}
