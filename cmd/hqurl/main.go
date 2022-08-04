package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hueristiq/url"
	"github.com/logrusorgru/aurora/v3"
)

type options struct {
	input   string
	fmtStr  string
	mode    string
	noColor bool
	unique  bool
	verbose bool
}

var (
	co options
)

func banner() {
	fmt.Fprintln(os.Stderr, aurora.BrightBlue(`
 _                      _ 
| |__   __ _ _   _ _ __| |
| '_ \ / _`+"`"+` | | | | '__| |
| | | | (_| | |_| | |  | |
|_| |_|\__, |\__,_|_|  |_| v1.0.0
          |_|             
`).Bold())
}

func init() {
	flag.BoolVar(&co.unique, "u", false, "")
	flag.BoolVar(&co.verbose, "v", false, "")
	flag.StringVar(&co.input, "input", "", "")
	flag.StringVar(&co.input, "i", "", "")

	flag.Usage = func() {
		banner()

		h := "USAGE:\n"
		h += "  hqurl [OPTIONS] [MODE] [FORMATSTRING]\n\n"

		h += "GENERAL OPTIONS:\n"
		h += "  -i                input file\n"
		h += "  -u                output unique values\n"
		h += "  -v                verbose mode: output URL parse errors\n\n"

		h += "MODE OPTIONS:\n"
		h += "  keys              keys from the query string (one per line)\n"
		h += "  values            values from the query string (one per line)\n"
		h += "  keypairs          `key=value` pairs from the query string (one per line)\n"
		h += "  hostnames         the hostname (e.g. sub.example.com)\n"
		h += "  paths             the request path (e.g. /users)\n"
		h += "  format            specify a custom format (see below)\n\n"

		h += "FORMAT DIRECTIVES:\n"
		h += "  %%                a literal percent character\n"
		h += "  %s                the request scheme (e.g. https)\n"
		h += "  %u                the user info (e.g. user:pass)\n"
		h += "  %h                the hostname (e.g. sub.example.com)\n"
		h += "  %S                the subdomain (e.g. sub)\n"
		h += "  %r                the root of domain (e.g. example)\n"
		h += "  %t                the TLD (e.g. com)\n"
		h += "  %P                the port (e.g. 8080)\n"
		h += "  %p                the path (e.g. /users)\n"
		h += "  %q                the raw query string (e.g. a=1&b=2)\n"
		h += "  %f                the page fragment (e.g. page-section)\n"
		h += "  %@                inserts an @ if user info is specified\n"
		h += "  %:                inserts a colon if a port is specified\n"
		h += "  %?                inserts a question mark if a query string exists\n"
		h += "  %#                inserts a hash if a fragment exists\n"
		h += "  %a                authority (alias for %u%@%d%:%P)\n\n"

		h += "EXAMPLES:\n"
		h += "  cat urls.txt | hqurl keys\n"
		h += "  cat urls.txt | hqurl format %s://%h%p?%q\n"

		fmt.Fprint(os.Stderr, h)
	}

	flag.Parse()

	co.mode = flag.Arg(0)
	co.fmtStr = flag.Arg(1)
}

// a urlProc is any function that accepts a URL and some
// kind of format string (which may not actually be used
// by some functions), and returns a slice of strings
// derived from that URL. It's not uncommon for a urlProc
// function to return a slice of length 1, but the return
// type remains a slice because *some* functions need to
// return multiple strings; e.g. the keys function.
type urlProc func(*url.URL, string) []string

// keys returns all of the keys used in the query string
// portion of the URL. E.g. for /?one=1&two=2&three=3 it
// will return []string{"one", "two", "three"}
func keys(u *url.URL, _ string) []string {
	out := make([]string, 0)
	for key := range u.Query() {
		out = append(out, key)
	}
	return out
}

// values returns all of the values in the query string
// portion of the URL. E.g. for /?one=1&two=2&three=3 it
// will return []string{"1", "2", "3"}
func values(u *url.URL, _ string) []string {
	out := make([]string, 0)
	for _, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, val)
		}
	}
	return out
}

// keyPairs returns all the key=value pairs in
// the query string portion of the URL. E.g for
// /?one=1&two=2&three=3 it will return
// []string{"one=1", "two=2", "three=3"}
func keyPairs(u *url.URL, _ string) []string {
	out := make([]string, 0)
	for key, vals := range u.Query() {
		for _, val := range vals {
			out = append(out, fmt.Sprintf("%s=%s", key, val))
		}
	}
	return out
}

// hostnames returns the domain portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"sub.example.com"}
func hostnames(u *url.URL, f string) []string {
	return format(u, "%h")
}

// domains returns the path portion of the URL. e.g.
// for http://sub.example.com/path it will return
// []string{"/path"}
func paths(u *url.URL, f string) []string {
	return format(u, "%p")
}

// format is a little bit like a special sprintf for
// URLs; it will return a single formatted string
// based on the URL and the format string. e.g. for
// http://example.com/path and format string "%d%p"
// it will return example.com/path
func format(u *url.URL, f string) []string {
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
			out.WriteRune('%')

		// the scheme; e.g. http
		case 's':
			out.WriteString(u.Scheme)

		// the userinfo; e.g. user:pass
		case 'u':
			if u.User != nil {
				out.WriteString(u.User.String())
			}

		// the domain; e.g. sub.example.com
		case 'h':
			out.WriteString(u.Hostname())

		// the port; e.g. 8080
		case 'P':
			out.WriteString(u.Port)

		// the subdomain; e.g. www
		case 'S':
			out.WriteString(u.SubDomain)

		// the root; e.g. example
		case 'r':
			out.WriteString(u.RootDomain)

		// the tld; e.g. com
		case 't':
			out.WriteString(u.TLD)

		// the path; e.g. /users
		case 'p':
			out.WriteString(u.EscapedPath())

		// the query string; e.g. one=1&two=2
		case 'q':
			out.WriteString(u.RawQuery)

		// the fragment / hash value; e.g. section-1
		case 'f':
			out.WriteString(u.Fragment)

		// an @ if user info is specified
		case '@':
			if u.User != nil {
				out.WriteRune('@')
			}

		// a colon if a port is specified
		case ':':
			if u.Port != "" {
				out.WriteRune(':')
			}

		// a question mark if there's a query string
		case '?':
			if u.RawQuery != "" {
				out.WriteRune('?')
			}

		// a hash if there is a fragment
		case '#':
			if u.Fragment != "" {
				out.WriteRune('#')
			}

		// the authority; e.g. user:pass@example.com:8080
		case 'a':
			out.WriteString(format(u, "%u%@%d%:%P")[0])

		// default to literal
		default:
			// output untouched
			out.WriteRune('%')
			out.WriteRune(r)
		}

		inFormat = false
	}

	return []string{out.String()}
}

func main() {
	procFn, ok := map[string]urlProc{
		"keys":      keys,
		"values":    values,
		"keypairs":  keyPairs,
		"hostnames": hostnames,
		"paths":     paths,
		"format":    format,
	}[co.mode]

	if !ok {
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", co.mode)
		return
	}

	seen := make(map[string]bool)

	var scanner *bufio.Scanner

	if co.input == "" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalln(errors.New("no stdin"))
		}

		if stat.Mode()&os.ModeNamedPipe == 0 {
			log.Fatalln(errors.New("no stdin"))
		}

		scanner = bufio.NewScanner(os.Stdin)
	} else {
		openedFile, err := os.Open(co.input)
		if err != nil {
			log.Fatalln(err)
		}
		defer openedFile.Close()

		scanner = bufio.NewScanner(openedFile)
	}

	for scanner.Scan() {
		parsedURL, err := url.Parse(url.Options{URL: scanner.Text()})
		if err != nil {
			if co.verbose {
				fmt.Fprintf(os.Stderr, "parse failure: %s\n", err)
			}

			continue
		}

		// some urlProc functions return multiple things,
		// so it's just easier to always get a slice and
		// loop over it instead of having two kinds of
		// urlProc functions.
		for _, val := range procFn(parsedURL, co.fmtStr) {

			// you do see empty values sometimes
			if val == "" {
				continue
			}

			if seen[val] && co.unique {
				continue
			}

			fmt.Println(val)

			// no point using up memory if we're outputting dupes
			if co.unique {
				seen[val] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}
}
