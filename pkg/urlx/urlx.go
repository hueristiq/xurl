package urlx

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type URL struct {
	// scheme, user (username & password), host (host or host:Port), path, query, fragment
	*url.URL
	Domain     string // e.g. sub.example.com
	ETLDPlus1  string // example.com
	SubDomain  string // e.g. sub
	RootDomain string // example
	TLD        string // com
	Port       string // 80
}

func Parse(rawURL string) (*URL, error) {
	return ParseWithDefaultScheme(rawURL, "http")
}

func ParseWithDefaultScheme(rawURL string, scheme string) (parsedURL *URL, err error) {
	rawURL = defaultScheme(rawURL, scheme)

	parsedURL = &URL{}

	parsedURL.URL, err = url.Parse(rawURL)
	if err != nil {
		return
	}

	for i := len(parsedURL.URL.Host) - 1; i >= 0; i-- {
		if parsedURL.URL.Host[i] == ':' {
			parsedURL.Domain = parsedURL.URL.Host[:i]
			parsedURL.Port = parsedURL.URL.Host[i+1:]
			break
		} else if parsedURL.URL.Host[i] < '0' || parsedURL.URL.Host[i] > '9' {
			parsedURL.Domain = parsedURL.URL.Host
		}
	}

	parsedURL.ETLDPlus1, err = publicsuffix.EffectiveTLDPlusOne(parsedURL.Domain)
	if err != nil {
		return
	}

	i := strings.Index(parsedURL.ETLDPlus1, ".")
	parsedURL.RootDomain = parsedURL.ETLDPlus1[0:i]
	parsedURL.TLD = parsedURL.ETLDPlus1[i+1:]

	if rest := strings.TrimSuffix(parsedURL.Domain, "."+parsedURL.ETLDPlus1); rest != parsedURL.Domain {
		parsedURL.SubDomain = rest
	}

	return
}
