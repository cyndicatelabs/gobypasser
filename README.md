# GoBypasser

A go utility to bypass URL 403 errors. Can easily be extended to include extra checks by modifying the `_bypassess.go` files.

## Usage

```
General Options:
    -f          A file containing a list of URLs to test. (e.g. https://google.com)
    -p          The base path you want to access
    -u          The host with schema (e.g. https://google.com)
    -v          Verbose output (default: false)

Attack Options:
    -headers    Cycle through all of the headers. (default: false)
    -paths      Cycle through all of the path bypasses. (default: false)
    -verbs      Cycle through all the verbs for the specified path. (default: false)

Filter Options:
    -hc         The response code(s) to hide (e.g. -hc 302 or -hc 404,400). (default: 403)
    -hs         The response size(s) to hide (e.g. -hs 4096 or -hs 4096,1024). (default: 0)
```

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

* Implement proper response size cannot use Response.Body
* Implement Stdin piping
