package main

import (
	"flag"
	"fmt"

	"github.com/cyndicatelabs/gobypasser/pkg/gobypasser"
)

func main() {

	o := gobypasser.Options{}
	gobypasser.SetDefaultOptions(&o)

	flag.StringVar(&o.FileOfUrls, "f", "", "A file containing a list of URLs to test. (e.g. https://google.com)")
	flag.StringVar(&o.BaseURL, "u", "", "The host with schema (e.g. https://google.com)")
	flag.StringVar(&o.BasePath, "p", "", "The base path you want to access")
	flag.StringVar(&o.PrimaryVerb, "X", "GET", "The HTTP verb to use for the baseline, path and header requests.")
	flag.BoolVar(&o.VerbBypasses, "verbs", false, "Cycle through all the verbs for the specified path.")
	flag.BoolVar(&o.PathBypasses, "paths", false, "Cycle through all of the path bypasses.")
	flag.BoolVar(&o.HeaderBypasses, "headers", false, "Cycle through all of the headers.")
	flag.StringVar(&o.UserAgent, "user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.72 Safari/537.36 Edg/88.0.705.81", "Specify the User-Agent for the requests.")
	flag.StringVar(&o.FilterResponseCode, "hc", "", "Additional response code(s) to hide (e.g. -hc 302 or -hc 404,400).")
	flag.StringVar(&o.FilterResponseSize, "hs", "", "Additional response size(s) to hide (e.g. -hs 4096 or -hs 4096,1024).")
	flag.BoolVar(&o.ShowAll, "show-all", false, "Show rows that match the baseline instead of hiding them.")
	flag.BoolVar(&o.ShowAll, "sa", false, "Alias for -show-all.")
	flag.BoolVar(&o.NoColor, "no-color", false, "Disable ANSI colors in output.")
	flag.BoolVar(&o.FollowRedirects, "L", false, "Follow redirects (default: show the raw 3xx response).")
	flag.BoolVar(&o.Verbose, "v", false, "Verbose output (show hidden rows plus per-request errors/timeouts).")
	flag.IntVar(&o.Threads, "t", 30, "Concurrent requests")

	flag.Usage = gobypasser.Usage
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	r, msg := gobypasser.VerifyOptions(&o)
	if !r {
		flag.Usage()
		fmt.Printf("[DEBUG] %s\n", msg)
		return
	}

	gobypasser.Banner()
	gobypasser.Start(&o)
}
