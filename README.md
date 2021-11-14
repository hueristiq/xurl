# urlx

![made with go](https://img.shields.io/badge/made%20with-Go-0040ff.svg) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0040ff.svg) [![open issues](https://img.shields.io/github/issues-raw/enenumxela/urlx.svg?style=flat&color=0040ff)](https://github.com/enenumxela/urlx/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/enenumxela/urlx.svg?style=flat&color=0040ff)](https://github.com/enenumxela/urlx/issues?q=is:issue+is:closed) [![license](https://img.shields.io/badge/License-MIT-gray.svg?colorB=0040FF)](https://github.com/enenumxela/urlx/blob/master/LICENSE) [![twitter](https://img.shields.io/badge/twitter-@enenumxela-0040ff.svg)](https://twitter.com/enenumxela)

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

You can download the pre-built binary for your platform from this repository's [releases](https://github.com/enenumxela/urlx/releases/) page, extract, then move it to your `$PATH`and you're ready to go.

#### From Source

urlx requires **go1.14+** to install successfully. Run the following command to get the repo

```bash
GO111MODULE=on go get -u -v github.com/enenumxela/urlx/cmd/urlx
```

#### From Github

```bash
git clone https://github.com/enenumxela/urlx.git && \
cd urlx/cmd/urlx/ && \
go build && \
mv urlx /usr/local/bin/ && \
urlx -h
```

## Usage

urlx works with URLs provided on stdin; they might come from a file like `urls.txt`:

```
$ cat example-urls.txt

https://sub.example.com/users?id=123&name=Sam
https://sub.example.com/orgs?org=ExCo#about
http://example.net/about#contact
```

You can extract:

* Hostnames from the URLs with the `hostnames` mode:

    ```
    $ cat urls.txt | urlx hostnames

    sub.example.com
    sub.example.com
    example.net
    ```

* Paths, with the `paths` mode:

    ```
    $ cat urls.txt | urlx paths

    /users
    /orgs
    /about
    ```

* Query String Keys, with the `keys` mode:

    ```
    $ cat urls.txt | urlx keys

    id
    name
    org
    ```

* Query String Values, with the `values` mode:

    ```
    $ cat urls.txt | urlx values

    123
    Sam
    ExCo
    ```

* Query String Key/Value Pairs , with the `keypairs` mode:

    ```
    $ cat urls.txt | urlx keypairs

    id=123
    name=Sam
    org=ExCo
    ```


* **NOTE:** You can use the `format` mode to specify a custom output format:

    ```
    $ cat urls.txt | urlx format %d%p

    sub.example.com/users
    sub.example.com/orgs
    example.net/about
    ```

    > For more format directives, checkout the help message `urlx -h` under `Format Directives`. 
    
    Any characters that don't match a format directive remain untouched:

    ```
    $ cat urls.txt | urlx -u format "%d (%s)"

    sub.example.com (https)
    example.net (http)
    ```

**Note** that if a URL does not include the data requested, there will be no output for that URL:

```
$ echo http://example.com | urlx format "%P"

$ echo http://example.com:8080 | urlx format "%P"
8080
```

To display help message for urlx use the `-h` flag:

```bash
$ urlx -h
```

help message:

```text
            _
 _   _ _ __| |_  __
| | | | '__| \ \/ /
| |_| | |  | |>  < 
 \__,_|_|  |_/_/\_\ v1.0.0

USAGE:
  urlx [OPTIONS] [MODE] [FORMATSTRING]

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
  cat urls.txt | urlx keys
  cat urls.txt | urlx format %s://%h%p?%q
```

## Credits

All credits to [Tom Hudson](https://github.com/tomnomnom), i took the initial code from his [unfurl](https://github.com/tomnomnom/unfurl).
## Contibution

[Issues](https://github.com/enenumxela/urlx/issues) and [Pull Requests](https://github.com/enenumxela/urlx/pulls) are welcome! 