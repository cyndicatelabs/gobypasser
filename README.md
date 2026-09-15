# GoBypasser

A go utility to bypass URL 403 errors. Can easily be extended to include extra checks by modifying the `_bypassess.go` files.

## Installation

* If you have go compiler installed: `go install github.com/cyndicatelabs/gobypasser@latest` - This will also update the tool.
* If not - `git clone https://github.com/cyndicatelabs/gobypasser ; cd gobypasser ; go build .`

## Usage

```
General Options:
    -u          The host with schema (e.g. https://google.com)
    -f          A file containing a list of URLs to test. (e.g. https://google.com)
    -p          The base path you want to access
    -X          The HTTP verb for the baseline, path and header requests. (default: GET)
    -t          Concurrent requests. (default: 30)
    -v          Verbose output — show hidden rows plus per-request errors. (default: false)
    -L          Follow redirects (default: show the raw 3xx response).
    -no-color   Disable ANSI colors in output.

Attack Options:
    -headers    Cycle through all of the headers. (default: false)
    -paths      Cycle through all of the path bypasses. (default: false)
    -verbs      Cycle through all the verbs for the specified path. (default: false)

Filter Options:
    -hc         Additional response code(s) to hide (e.g. -hc 302 or -hc 404,400).
    -hs         Additional response size(s) to hide (e.g. -hs 4096 or -hs 4096,1024).
    -show-all   Show rows that match the baseline instead of hiding them. (alias: -sa)
```

## How it works

Before probing each target, GoBypasser sends one unmodified request
(`{-X verb} {target}/{-p path}`) and records its status and size as the
**baseline** — normally the 403 you're trying to defeat. Every bypass is then
compared against that baseline, and rows matching it (same status *and* size) are
hidden as noise, so a single successful bypass no longer drowns in hundreds of
identical 403s. Pass `-show-all` (or `-v`) to see the suppressed rows.

Header bypasses are split by technique: IP-spoof and host-override headers
(`X-Forwarded-For: 127.0.0.1`, …) are sent against the protected path, while
URL-rewrite headers (`X-Original-URL`, `X-Rewrite-URL`, …) carry the path and are
sent against the site root.

## Example Output

```
$ gobypasser -u http://127.0.0.1:8080 -p admin -paths -hc 404

                GoBypasser 0.1 - https://www.github.com/cyndicatelabs/gobypasser -

            A tool to aid finding URL bypasses using a number of different techniques.
            
Settings:
  Base URL   : http://127.0.0.1:8080
  Base Path  : /admin
  User-Agent : Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.72 Safari/537.36 Edg/88.0.705.81

Response Code   Response Size   Verb                 Path                                                                                       Custom Header
__________________________________________________________________________________________________________________________________________________________________________
200             0               GET                  //admin/../                                                                                N/A
```

## TODO

* Implement Stdin piping
