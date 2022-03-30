package urlx_test

import (
	"testing"

	"github.com/enenumxela/urlx/pkg/urlx"
)

func TestDefaultScheme(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{in: "localhost", out: "http://localhost"},
		{in: "example.com", out: "http://example.com"},
		{in: "https://example.com", out: "https://example.com"},
		{in: "://example.com", out: "http://example.com"},
		{in: "//example.com", out: "http://example.com"},
	}

	for _, tt := range tests {
		URL := urlx.DefaultScheme(tt.in)

		if URL != tt.out {
			t.Errorf(`"%s": got "%s", want "%v"`, tt.in, URL, tt.out)
		}
	}
}
