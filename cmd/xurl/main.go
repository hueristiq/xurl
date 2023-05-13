package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"sync"

	hqlog "github.com/hueristiq/hqgoutils/log"
	"github.com/hueristiq/hqgoutils/log/formatter"
	"github.com/hueristiq/hqgoutils/log/levels"
	hqurl "github.com/hueristiq/hqgoutils/url"
	"github.com/hueristiq/xurl/internal/configuration"
	"github.com/hueristiq/xurl/internal/processor"
	"github.com/spf13/pflag"
)

var (
	input string

	monochrome bool
	unique     bool
	verbosity  string

	mode   string
	fmtStr string
)

func init() {
	pflag.StringVarP(&input, "input", "i", "", "")
	pflag.BoolVarP(&unique, "unique", "u", false, "")
	pflag.BoolVarP(&monochrome, "monochrome", "m", false, "")
	pflag.StringVarP(&verbosity, "verbosity", "v", string(levels.LevelInfo), "")

	pflag.CommandLine.SortFlags = false
	pflag.Usage = func() {
		fmt.Fprintln(os.Stderr, configuration.BANNER)

		h := "\nUSAGE:\n"
		h += "  xurl [MODE] [FORMATSTRING] [OPTIONS]\n"

		h += "\nINPUT:\n"
		h += "  -i, --input       input file (use `-` to get from stdin)\n"

		h += "\nOUTPUT:\n"
		h += "  -m, --monochrome  disable output content coloring\n"
		h += "  -u, --unique      output unique values\n"
		h += fmt.Sprintf("  -v, --verbosity   debug, info, warning, error, fatal or silent (default: %s)\n\n", string(levels.LevelInfo))

		h += "\nMODE:\n"
		h += "  domains           the hostname (e.g. sub.example.com)\n"
		h += "  apexes            the apex domain (e.g. example.com from sub.example.com)\n"
		h += "  paths             the request path (e.g. /users)\n"
		h += "  query             `key=value` pairs from the query string (one per line)\n"
		h += "  params            keys from the query string (one per line)\n"
		h += "  values            values from the query string (one per line)\n"
		h += "  format            custom format (see below)\n\n"

		h += "FORMAT DIRECTIVES:\n"
		h += "  %%                a literal percent character\n"
		h += "  %s                the request scheme (e.g. https)\n"
		h += "  %u                the user info (e.g. user:pass)\n"
		h += "  %d                the domain (e.g. sub.example.com)\n"
		h += "  %S                the subdomain (e.g. sub)\n"
		h += "  %r                the root of domain (e.g. example)\n"
		h += "  %t                the TLD (e.g. com)\n"
		h += "  %P                the port (e.g. 8080)\n"
		h += "  %p                the path (e.g. /users)\n"
		h += "  %e                the path's file extension (e.g. jpg, html)\n"
		h += "  %q                the raw query string (e.g. a=1&b=2)\n"
		h += "  %f                the page fragment (e.g. page-section)\n"
		h += "  %@                inserts an @ if user info is specified\n"
		h += "  %:                inserts a colon if a port is specified\n"
		h += "  %?                inserts a question mark if a query string exists\n"
		h += "  %#                inserts a hash if a fragment exists\n"
		h += "  %a                authority (alias for %u%@%d%:%P)\n\n"

		h += "EXAMPLES:\n"
		h += "  cat urls.txt | xurl params -i -\n"
		h += "  cat urls.txt | xurl format %s://%h%p?%q -i -\n"

		fmt.Fprint(os.Stderr, h)
	}

	pflag.Parse()

	mode = pflag.Arg(0)
	fmtStr = pflag.Arg(1)

	hqlog.DefaultLogger.SetMaxLevel(levels.LevelStr(verbosity))
	hqlog.DefaultLogger.SetFormatter(formatter.NewCLI(&formatter.CLIOptions{
		Colorize: !monochrome,
	}))
}

func main() {
	// mode
	procFn, ok := map[string]processor.Extractor{
		"domains": processor.Domains,
		"apexes":  processor.Apexes,
		"paths":   processor.Paths,
		"query":   processor.Query,
		"params":  processor.Parameters,
		"values":  processor.Values,
		"format":  processor.Format,
	}[mode]

	if !ok {
		hqlog.Fatal().Msgf("unknown mode: %s", mode)
	}

	// input URLs
	URLs := make(chan string)

	go func() {
		defer close(URLs)

		var (
			err  error
			file *os.File
			stat fs.FileInfo
		)

		switch {
		case input != "" && input == "-":
			stat, err = os.Stdin.Stat()
			if err != nil {
				hqlog.Fatal().Msg("no stdin")
			}

			if stat.Mode()&os.ModeNamedPipe == 0 {
				hqlog.Fatal().Msg("no stdin")
			}

			file = os.Stdin
		case input != "" && input != "-":
			file, err = os.Open(input)
			if err != nil {
				hqlog.Fatal().Msg(err.Error())
			}
		default:
			hqlog.Fatal().Msg("xurl takes input from stdin or file using ")
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			URL := scanner.Text()

			if URL != "" {
				URLs <- URL
			}
		}

		if err := scanner.Err(); err != nil {
			hqlog.Fatal().Msg(err.Error())
		}
	}()

	// process URLs
	wg := &sync.WaitGroup{}

	seen := &sync.Map{}

	for URL := range URLs {
		wg.Add(1)

		go func(URL string) {
			defer wg.Done()

			parsedURL, err := hqurl.Parse(URL)
			if err != nil {
				hqlog.Error().Msgf("parse failure: %s", err.Error())

				return
			}

			for _, value := range procFn(parsedURL, fmtStr) {
				if value == "" {
					continue
				}

				if unique {
					_, loaded := seen.LoadOrStore(value, struct{}{})
					if loaded {
						continue
					}
				}

				hqlog.Print().Msg(value)
			}
		}(URL)
	}

	wg.Wait()
}
