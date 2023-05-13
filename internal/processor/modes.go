package processor

import (
	"bytes"
	"fmt"
	"strings"

	hqurl "github.com/hueristiq/hqgoutils/url"
)

// format is a little bit like a special sprintf for
// URLs; it will return a single formatted string
// based on the URL and the format string. e.g. for
// http://example.com/path and format string "%d%p"
// it will return example.com/path
func Format(u *hqurl.URL, f string) []string {
	out := &bytes.Buffer{}

	inFormat := false

	for _, r := range f {
		if r == '%' && !inFormat {
			inFormat = true

			continue
		}

		if !inFormat {
			out.WriteRune(r)

			continue
		}

		switch r {
		// a literal percent rune
		case '%':
			out.WriteByte('%')

		// the scheme; e.g. http
		case 's':
			out.WriteString(u.Scheme)

		// the userinfo; e.g. user:pass
		case 'u':
			if u.User != nil {
				out.WriteString(u.User.String())
			}

		// the domain; e.g. sub.example.com
		case 'd':
			out.WriteString(u.Hostname())

		// the port; e.g. 8080
		case 'P':
			out.WriteString(u.Port)

		// the subdomain; e.g. www
		case 'S':
			out.WriteString(u.Subdomain)

		// the root; e.g. example
		case 'r':
			out.WriteString(u.RootDomain)

		// the tld; e.g. com
		case 't':
			out.WriteString(u.TLD)

		// the path; e.g. /users
		case 'p':
			out.WriteString(u.EscapedPath())

		// the paths's file extension
		case 'e':
			paths := strings.Split(u.EscapedPath(), "/")
			if len(paths) > 1 {
				parts := strings.Split(paths[len(paths)-1], ".")
				if len(parts) > 1 {
					out.WriteString(parts[len(parts)-1])
				}
			} else {
				parts := strings.Split(u.EscapedPath(), ".")
				if len(parts) > 1 {
					out.WriteString(parts[len(parts)-1])
				}
			}

		// the query string; e.g. one=1&two=2
		case 'q':
			out.WriteString(u.RawQuery)

		// the fragment / hash value; e.g. section-1
		case 'f':
			out.WriteString(u.Fragment)

		// an @ if user info is specified
		case '@':
			if u.User != nil {
				out.WriteByte('@')
			}

		// a colon if a port is specified
		case ':':
			if u.Port != "" {
				out.WriteByte(':')
			}

		// a question mark if there's a query string
		case '?':
			if u.RawQuery != "" {
				out.WriteByte('?')
			}

		// a hash if there is a fragment
		case '#':
			if u.Fragment != "" {
				out.WriteByte('#')
			}

		// the authority; e.g. user:pass@example.com:8080
		case 'a':
			out.WriteString(Format(u, "%u%@%d%:%P")[0])

		// default to literal
		default:
			// output untouched
			out.WriteByte('%')
			out.WriteRune(r)
		}

		inFormat = false
	}

	return []string{out.String()}
}

// hostnames returns the domain portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"sub.example.com"}
func Domains(u *hqurl.URL, _ string) []string {
	return Format(u, "%d")
}

// Apexes return the apex portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"example.com"}
func Apexes(u *hqurl.URL, _ string) []string {
	return Format(u, "%r.%t")
}

// domains returns the path portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"/path"}
func Paths(u *hqurl.URL, _ string) []string {
	return Format(u, "%p")
}

// keyPairs returns all the key=value pairs in
// the query string portion of the URL. E.g for
// /?one=1&two=2&three=3 it will return
// []string{"one=1", "two=2", "three=3"}
func Query(u *hqurl.URL, _ string) []string {
	out := make([]string, 0)

	// param:value
	for key, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, fmt.Sprintf("%s=%s", key, val))
		}
	}

	return out
}

// Parameters returns all of the keys used in the query string
// portion of the URL. E.g. for /?one=1&two=2&three=3 it
// will return []string{"one", "two", "three"}
func Parameters(u *hqurl.URL, _ string) []string {
	out := make([]string, 0)

	// param:value
	for key := range u.Query() {
		out = append(out, key)
	}

	return out
}

// values returns all of the values in the query string
// portion of the URL. E.g. for /?one=1&two=2&three=3 it
// will return []string{"1", "2", "3"}
func Values(u *hqurl.URL, _ string) []string {
	out := make([]string, 0)

	// param:value
	for _, value := range u.Query() {
		// value: [string]{items...}
		out = append(out, value...)
	}

	return out
}
