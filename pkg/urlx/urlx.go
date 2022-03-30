package urlx

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// URL represents a parsed URL (technically, a URI reference).
// https://sub.example.com:8080
type URL struct {
	// Scheme      string    // e.g https
	// Opaque      string    // encoded opaque data
	// User        *Userinfo // username and password information
	// Host        string    // e.g. sub.example.com, sub.example.com:8080
	// Path        string    // path (relative paths may omit leading slash)
	// RawPath     string    // encoded path hint (see EscapedPath method)
	// ForceQuery  bool      // append a query ('?') even if RawQuery is empty
	// RawQuery    string    // encoded query values, without '?'
	// Fragment    string    // fragment for references, without '#'
	// RawFragment string    // encoded fragment hint (see EscapedFragment method)
	*url.URL
	Domain      string // e.g. sub.example.com
	ETLDPlusOne string // e.g. example.com
	SubDomain   string // e.g. sub
	RootDomain  string // e.g. example
	TLD         string // e.g. com
	Port        string // e.g. 8080
}

// Parse parses a raw url into a urlx.URL structure.
//
// It uses the url.Parse() internally, but it slightly changes
// its behavior:
// 1. It forces the default scheme and port to http
// 2. It favors absolute paths over relative ones, thus "example.com"
//    is parsed into url.Host instead of url.Path.
// 3. It lowercases the Host (not only the Scheme).
func Parse(rawURL string) (*URL, error) {
	return ParseWithDefaultScheme(rawURL)
}

func ParseWithDefaultScheme(rawURL string) (parsedURL *URL, err error) {
	rawURL = DefaultScheme(rawURL)

	parsedURL = &URL{}

	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		return
	}

	// Domain + Port
	for i := len(parsedURL.URL.Host) - 1; i >= 0; i-- {
		if parsedURL.URL.Host[i] == ':' {
			parsedURL.Domain = parsedURL.URL.Host[:i]
			parsedURL.Port = parsedURL.URL.Host[i+1:]
			break
		} else if parsedURL.URL.Host[i] < '0' || parsedURL.URL.Host[i] > '9' {
			parsedURL.Domain = parsedURL.URL.Host
		}
	}

	// ETLDPlusOne
	parsedURL.ETLDPlusOne, err = publicsuffix.EffectiveTLDPlusOne(parsedURL.Domain)
	if err != nil {
		return
	}

	// RootDomain + TLD
	i := strings.Index(parsedURL.ETLDPlusOne, ".")
	parsedURL.RootDomain = parsedURL.ETLDPlusOne[0:i]
	parsedURL.TLD = parsedURL.ETLDPlusOne[i+1:]

	// Subdomain
	if rest := strings.TrimSuffix(parsedURL.Domain, "."+parsedURL.ETLDPlusOne); rest != parsedURL.Domain {
		parsedURL.SubDomain = rest
	}

	return
}

// DefaultScheme forces default scheme to `http` scheme, so net/url.Parse() doesn't
// put both host and path into the (relative) path.
func DefaultScheme(URL string) (URLWithScheme string) {
	URLWithScheme = URL

	if strings.Index(URLWithScheme, "//") == 0 {
		URLWithScheme = "http:" + URLWithScheme
	}

	if strings.Contains(URLWithScheme, "://") && !strings.HasPrefix(URLWithScheme, "http") {
		URLWithScheme = "http" + URLWithScheme
	}

	if !strings.Contains(URLWithScheme, "://") {
		URLWithScheme = "http://" + URLWithScheme
	}

	return
}
