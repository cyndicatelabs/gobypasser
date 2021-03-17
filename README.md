# GoBypasser

A simple tool to bypass URL path restrictions. Can easily be extended to include extra checks by modifying the `_bypassess.go` files.

## Usage

```
                GoBypasser 0.1 - https://www.github.com/cyndicatelabs/gobypasser -
                                Credits: @_g0dmode, @dcocking7

            A tool to aid finding URL bypasses using a number of different techniques.

General Options:
    -p          The base path you want to access
    -u          The host with schema (e.g. https://google.com)
    -v          Verbose output (default: false)

Attack Options:
    -headers    Cycle through all of the headers. (default: false)
    -paths      Cycle through all of the path bypasses. (default: false)
    -verbs      Cycle through all the verbs for the specified path. (default: false)

Filter Options:
    -hc         The response code(s) to hide (e.g. --hc 302 or --hc 404,400). (default: 403)
    -hs         The response size(s) to hide (e.g. --hs 4096 or --hs 4096,1024). (default: 0)

EXAMPLE USAGE:
        gobypasser -u https://www.google.com/ -p /api
        gobypasser -u https://www.google.com/ -p /api --hc 302,404 --verbs
        gobypasser -u https://www.google.com/ -p /api --hs 1024 --headers

```

## Example Output

```
[21:14:54 - Mon Mar 15] [eth0:172.24.60.50/20] [mitch@DESKTOP-XXXXXX]
[/gobypasser] $ ./gobypasser -u http://127.0.0.1:8080 -p admin -paths --hc 404

                GoBypasser 0.1 - https://www.github.com/cyndicatelabs/gobypasser -
                                Credits: @_g0dmode, @dcocking7

            A tool to aid finding URL bypasses using a number of different techniques.
            
Settings:
  Base URL   : http://127.0.0.1:8080
  Base Path  : /admin
  User-Agent : Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.72 Safari/537.36 Edg/88.0.705.81

Response Code   Response Size   Verb                 Path                                                                                       Custom Header
__________________________________________________________________________________________________________________________________________________________________________
200             0               GET                  //admin/../                                                                                N/A
```

## Credits

* Mitchell Hines - https://twitter.com/@_g0dmode
* Daniel Cocking - https://twitter.com/@dcocking7

## TODO

* Implement proper response size cannot use Response.Body
