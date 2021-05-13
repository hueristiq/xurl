package urlx

import (
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type URL struct {
	// scheme, user (username & password), host (host or host:Port), path, query, fragment
	*url.URL
	// Host (e.g. sub.example.com), RootDomain (e.g. example), SubDOmain (e.g. sub), TLD (e.g. com), Port (e.g. 80)
	Domain, RootDomain, SubDomain, TLD, Port string
}

func Parse(rawURL string) (parsedURL *URL, err error) {
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

	etldPlus1, err := publicsuffix.EffectiveTLDPlusOne(parsedURL.Domain)
	if err != nil {
		return nil, err
	}

	i := strings.Index(etldPlus1, ".")
	parsedURL.RootDomain = etldPlus1[0:i]
	parsedURL.TLD = etldPlus1[i+1:]

	if rest := strings.TrimSuffix(parsedURL.Domain, "."+etldPlus1); rest != parsedURL.Domain {
		parsedURL.SubDomain = rest
	}

	return
}
