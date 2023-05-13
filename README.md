# xurl

![made with go](https://img.shields.io/badge/made%20with-Go-0000FF.svg) [![release](https://img.shields.io/github/release/hueristiq/xurl?style=flat&color=0000FF)](https://github.com/hueristiq/xurl/releases) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=0000FF)](https://github.com/hueristiq/xurl/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0000FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/xurl.svg?style=flat&color=0000FF)](https://github.com/hueristiq/xurl/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/xurl.svg?style=flat&color=0000FF)](https://github.com/hueristiq/xurl/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-0000FF.svg)](https://github.com/hueristiq/xurl/blob/master/CONTRIBUTING.md)

`xurl` is a command-line interface (CLI) utility pull out bits of URLs.

## Resources

* [Installation](#installation)
	* [Install release binaries](#install-release-binaries)
	* [Install source](#install-sources)
		* [`go install ...`](#go-install)
		* [`go build ...` the development Version](#go-build--the-development-version)
* [Usage](#usage)
	* [Examples](#examples)
		* [Domains](#domains)
		* [Apex Domains](#apex-domains)
		* [Paths](#paths)
		* [Query String Key/Value Pairs](#query-string-keyvalue-pairs)
		* [Query String Keys (Parameters)](#query-string-keys-parameters)
		* [Query String Values](#query-string-values)
		* [Custom Formats](#custom-formats)
* [Credits](#credits)
* [Contribution](#contribution)
* [Licensing](#licensing)

## Installation

### Install release binaries

Visit the [releases page](https://github.com/hueristiq/xurl/releases) and find the appropriate archive for your operating system and architecture. Download the archive from your browser or copy its URL and retrieve it with `wget` or `curl`:

* ...with `wget`:

	```bash
	wget https://github.com/hueristiq/xurl/releases/download/v<version>/xurl-<version>-linux-amd64.tar.gz
	```

* ...or, with `curl`:

	```bash
	curl -OL https://github.com/hueristiq/xurl/releases/download/v<version>/xurl-<version>-linux-amd64.tar.gz
	```

...then, extract the binary:

```bash
tar xf xurl-<version>-linux-amd64.tar.gz
```

> **TIP:** The above steps, download and extract, can be combined into a single step with this onliner
> 
> ```bash
> curl -sL https://github.com/hueristiq/xurl/releases/download/v<version>/xurl-<version>-linux-amd64.tar.gz | tar -xzv
> ```

**NOTE:** On Windows systems, you should be able to double-click the zip archive to extract the `xurl` executable.

...move the `xurl` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

```bash
sudo mv xurl /usr/local/bin/
```

**NOTE:** Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xurl` to their `PATH`.

### Install source

Before you install from source, you need to make sure that Go is installed on your system. You can install Go by following the official instructions for your operating system. For this, we will assume that Go is already installed.

#### `go install ...`

```bash
go install -v github.com/hueristiq/xurl/cmd/xurl@latest
```

#### `go build ...` the development Version

* Clone the repository

	```bash
	git clone https://github.com/hueristiq/xurl.git 
	```

* Build the utility

	```bash
	cd xurl/cmd/xurl && \
	go build .
	```

* Move the `xurl` binary to somewhere in your `PATH`. For example, on GNU/Linux and OS X systems:

	```bash
	sudo mv xurl /usr/local/bin/
	```

	**NOTE:** Windows users can follow [How to: Add Tool Locations to the PATH Environment Variable](https://msdn.microsoft.com/en-us/library/office/ee537574(v=office.14).aspx) in order to add `xurl` to their `PATH`.


**NOTE:** While the development version is a good way to take a peek at `xurl`'s latest features before they get released, be aware that it may have bugs. Officially released versions will generally be more stable.

## Usage

To display help message for xurl use the `-h` flag:

```bash
$ xurl -h
```

help message:

```text
                 _ 
__  ___   _ _ __| |
\ \/ / | | | '__| |
 >  <| |_| | |  | |
/_/\_\\__,_|_|  |_| v0.0.0

A CLI utility to pull out bits of URLs.

USAGE:
  xurl [MODE] [FORMATSTRING] [OPTIONS]

INPUT:
  -i, --input       input file (use `-` to get from stdin)

OUTPUT:
  -m, --monochrome  disable output content coloring
  -u, --unique      output unique values
  -v, --verbosity   debug, info, warning, error, fatal or silent (default: info)

MODE:
  domains           the hostname (e.g. sub.example.com)
  apexes            the apex domain (e.g. example.com from sub.example.com)
  paths             the request path (e.g. /users)
  query             `key=value` pairs from the query string (one per line)
  params            keys from the query string (one per line)
  values            values from the query string (one per line)
  format            custom format (see below)

FORMAT DIRECTIVES:
  %%                a literal percent character
  %s                the request scheme (e.g. https)
  %u                the user info (e.g. user:pass)
  %d                the domain (e.g. sub.example.com)
  %S                the subdomain (e.g. sub)
  %r                the root of domain (e.g. example)
  %t                the TLD (e.g. com)
  %P                the port (e.g. 8080)
  %p                the path (e.g. /users)
  %e                the path's file extension (e.g. jpg, html)
  %q                the raw query string (e.g. a=1&b=2)
  %f                the page fragment (e.g. page-section)
  %@                inserts an @ if user info is specified
  %:                inserts a colon if a port is specified
  %?                inserts a question mark if a query string exists
  %#                inserts a hash if a fragment exists
  %a                authority (alias for %u%@%d%:%P)

EXAMPLES:
  cat urls.txt | xurl params -i -
  cat urls.txt | xurl format %s://%h%p?%q -i -
```

### Examples

```
$ cat urls.txt

https://sub.example.com/users?id=123&name=Sam
https://sub.example.com/orgs?org=ExCo#about
http://example.net/about#contact
```

#### Domains

You can extract the domains from the URLs with the `domains` mode:

```
$ cat urls.txt | xurl domains -i -

sub.example.com
sub.example.com
example.net
```

If you don't want to output duplicate values you can use the `-u` or `--unique` flag:

	```
	$ cat urls.txt | xurl domains  -i - --unique
	sub.example.com
	example.net
	```

The `-u`/`--unique` flag works for all modes.

#### Apex Domains

You can extract the apex part of the domain (e.g. the `example.com` in `http://sub.example.com`) using the `apexes` mode:

```
$ cat urls.txt | unfurl apexes -i - -u
example.com
example.net
```

#### Paths

```
$ cat urls.txt | xurl paths -i -

/users
/orgs
/about
```

#### Query String Key/Value Pairs

```
$ cat urls.txt | xurl query -i -

id=123
name=Sam
org=ExCo
```

#### Query String Keys (Parameters)

```
$ cat urls.txt | xurl params -i -

id
name
org
```

#### Query String Values

```
$ cat urls.txt | xurl values -i -

123
Sam
ExCo
```

#### Custom Formats

You can use the `format` mode to specify a custom output format:

```
$ cat urls.txt | xurl format %d%p -i -

sub.example.com/users
sub.example.com/orgs
example.net/about
```

The available format directives are:

```
%%  A literal percent character
%s  The request scheme (e.g. https)
%u  The user info (e.g. user:pass)
%d  The domain (e.g. sub.example.com)
%S  The subdomain (e.g. sub)
%r  The root of domain (e.g. example)
%t  The TLD (e.g. com)
%P  The port (e.g. 8080)
%p  The path (e.g. /users)
%e  The path's file extension (e.g. jpg, html)
%q  The raw query string (e.g. a=1&b=2)
%f  The page fragment (e.g. page-section)
%@  Inserts an @ if user info is specified
%:  Inserts a colon if a port is specified
%?  Inserts a question mark if a query string exists
%#  Inserts a hash if a fragment exists
%a  Authority (alias for %u%@%d%:%P)
```

> For more format directives, checkout the help message `xurl -h` under `Format Directives`. 

Any characters that don't match a format directive remain untouched:

```
$ cat urls.txt | xurl format "%d (%s)"  -i - -u

sub.example.com (https)
example.net (http)
```

**Note** that if a URL does not include the data requested, there will be no output for that URL:

```
$ echo http://example.com | xurl format "%P"  -i -

$ echo http://example.com:8080 | xurl format "%P" -i -
8080
```

## Credits

* [Tom Hudson](https://github.com/tomnomnom), we took the initial code from his [unfurl](https://github.com/tomnomnom/unfurl).
## Contibution

[Issues](https://github.com/hueristiq/xurl/issues) and [Pull Requests](https://github.com/hueristiq/xurl/pulls) are welcome! Check out the [contribution guidelines.](./CONTRIBUTING.md)
## Licensing

This utility is distributed under the [MIT license](./LICENSE)