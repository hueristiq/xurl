# xurl

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/xurl.svg?style=flat&color=0040ff)](https://github.com/hueristiq/xurl/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/xurl.svg?style=flat&color=0040ff)](https://github.com/hueristiq/xurl/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/License-MIT-gray.svg?colorB=0040FF)](https://github.com/hueristiq/xurl/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@itshueristiq-0040ff.svg)](https://twitter.com/itshueristiq)

A [go(golang)](http://golang.org/) utility for URLs parsing & pull out bits of the URLS.

## Resources

* [Installation](#installation)
	* [From Binary](#from-binary)
	* [From source](#from-source)
	* [From github](#from-github)
* [Usage](#usage)
* [Credits](#credits)
* [Contribution](#contribution)

## Installation

#### From Binary

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/hueristiq/xurl/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

xurl requires **go1.17+** to install successfully. Run the following command to get the repo

```bash
go install -v github.com/hueristiq/xurl/cmd/xurl
```

#### From Github

```bash
git clone https://github.com/hueristiq/xurl.git && \
cd xurl/cmd/xurl/ && \
go build && \
mv xurl /usr/local/bin/ && \
xurl -h
```

## Usage

xurl works with URLs provided on stdin; they might come from a file like `urls.txt`:

```
$ cat example-urls.txt

https://sub.example.com/users?id=123&name=Sam
https://sub.example.com/orgs?org=ExCo#about
http://example.net/about#contact
```

You can extract:

* Hostnames from the URLs with the `hostnames` mode:

    ```
    $ cat urls.txt | xurl hostnames

    sub.example.com
    sub.example.com
    example.net
    ```

* Paths, with the `paths` mode:

    ```
    $ cat urls.txt | xurl paths

    /users
    /orgs
    /about
    ```

* Query String Keys, with the `keys` mode:

    ```
    $ cat urls.txt | xurl keys

    id
    name
    org
    ```

* Query String Values, with the `values` mode:

    ```
    $ cat urls.txt | xurl values

    123
    Sam
    ExCo
    ```

* Query String Key/Value Pairs , with the `keypairs` mode:

    ```
    $ cat urls.txt | xurl keypairs

    id=123
    name=Sam
    org=ExCo
    ```


* **NOTE:** You can use the `format` mode to specify a custom output format:

    ```
    $ cat urls.txt | xurl format %d%p

    sub.example.com/users
    sub.example.com/orgs
    example.net/about
    ```

    > For more format directives, checkout the help message `xurl -h` under `Format Directives`. 
    
    Any characters that don't match a format directive remain untouched:

    ```
    $ cat urls.txt | xurl -u format "%d (%s)"

    sub.example.com (https)
    example.net (http)
    ```

**Note** that if a URL does not include the data requested, there will be no output for that URL:

```
$ echo http://example.com | xurl format "%P"

$ echo http://example.com:8080 | xurl format "%P"
8080
```

To display help message for xurl use the `-h` flag:

```bash
$ xurl -h
```

help message:

```text
            _
 _   _ _ __| |_  __
| | | | '__| \ \/ /
| |_| | |  | |>  < 
 \__,_|_|  |_/_/\_\ v1.0.0

USAGE:
  xurl [OPTIONS] [MODE] [FORMATSTRING]

GENERAL OPTIONS:
  -i                input file
  -u                output unique values
  -v                verbose mode: output URL parse errors

MODE OPTIONS:
  keys              keys from the query string (one per line)
  values            values from the query string (one per line)
  keypairs          `key=value` pairs from the query string (one per line)
  hostnames         the hostname (e.g. sub.example.com)
  paths             the request path (e.g. /users)
  format            specify a custom format (see below)

FORMAT DIRECTIVES:
  %%                a literal percent character
  %s                the request scheme (e.g. https)
  %u                the user info (e.g. user:pass)
  %h                the hostname (e.g. sub.example.com)
  %S                the subdomain (e.g. sub)
  %r                the root of domain (e.g. example)
  %t                the TLD (e.g. com)
  %P                the port (e.g. 8080)
  %p                the path (e.g. /users)
  %q                the raw query string (e.g. a=1&b=2)
  %f                the page fragment (e.g. page-section)
  %@                inserts an @ if user info is specified
  %:                inserts a colon if a port is specified
  %?                inserts a question mark if a query string exists
  %#                inserts a hash if a fragment exists
  %a                authority (alias for %u%@%d%:%P)

EXAMPLES:
  cat urls.txt | xurl keys
  cat urls.txt | xurl format %s://%h%p?%q
```

## Credits

All credits to [Tom Hudson](https://github.com/tomnomnom), i took the initial code from his [unfurl](https://github.com/tomnomnom/unfurl).
## Contibution

[Issues](https://github.com/hueristiq/xurl/issues) and [Pull Requests](https://github.com/hueristiq/xurl/pulls) are welcome! 