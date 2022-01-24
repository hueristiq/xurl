package urlx

import "strings"

func defaultScheme(rawURL, scheme string) string {
	// Force default http scheme, so net/url.Parse() doesn't
	// put both host and path into the (relative) path.
	if strings.Index(rawURL, "//") == 0 {
		// Leading double slashes (any scheme). Force http.
		rawURL = scheme + ":" + rawURL
	}
	if !strings.Contains(rawURL, "://") {
		// Missing scheme. Force http.
		rawURL = scheme + "://" + rawURL
	}
	return rawURL
}
