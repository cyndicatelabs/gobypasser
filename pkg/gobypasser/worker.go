package gobypasser

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// titleCase upper-cases the first letter of each slash-separated segment.
func titleCase(path string) string {
	segs := strings.Split(path, "/")
	for i, s := range segs {
		if s != "" {
			segs[i] = strings.ToUpper(s[:1]) + s[1:]
		}
	}
	return strings.Join(segs, "/")
}

// substitute expands the template placeholders for a given target.
func substitute(tmpl, url, basePath string) string {
	r := strings.NewReplacer(
		"{base_url}", url,
		"{base_path_upper}", strings.ToUpper(basePath),
		"{base_path_title}", titleCase(basePath),
		"{base_path}", basePath,
	)
	return r.Replace(tmpl)
}

// computeBaseline issues one unmodified request for the target and records its
// status/size. On error it returns an un-OK baseline (comparison disabled).
func computeBaseline(client HttpClient, url, verb, basePath string) Baseline {
	req := NewHttpRequest(client, fmt.Sprintf("%s/%s", url, basePath), verb)
	r := MakeHttpRequest(client, req)
	if r.Err != nil {
		return Baseline{OK: false}
	}
	return Baseline{StatusCode: r.StatusCode, Size: r.Size, OK: true}
}

func Start(o *Options) {

	start := time.Now()

	fmt.Println("Settings:")
	if o.BaseURL != "" {
		fmt.Printf("  Base URL   : %s\n", o.BaseURL)
	} else if o.FileOfUrls != "" {
		fmt.Printf("  Number of URLs : %d\n", len(o.UrlList))
	}
	fmt.Printf("  Base Path  : %s\n", o.BasePath)
	fmt.Printf("  Verb       : %s\n", o.PrimaryVerb)
	fmt.Printf("  User-Agent : %s\n", o.UserAgent)
	fmt.Println()

	var wg sync.WaitGroup
	MyClient := NewHttpClient(o)
	headerBypasses := BuildHeaderBypasses()

	PrintTableHeader()

	buffer := len(o.UrlList) * (len(VerbBypasses) + len(headerBypasses) + len(PathBypasses) + 1)
	if buffer < 1 {
		buffer = 1
	}
	results := make(chan *Result, buffer)
	workerPool := make(chan struct{}, o.Threads)

	// run wraps a single request: acquire a slot, build the result, tag it
	// against the baseline, and push it onto the results channel.
	run := func(baseline Baseline, url, method, customHeader, headerValue string) {
		wg.Add(1)
		workerPool <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-workerPool }()
			req := NewHttpRequest(MyClient, url, method)
			if customHeader != "" {
				req.Header.Set(customHeader, headerValue)
			}
			result := MakeHttpRequest(MyClient, req)
			if result.Err == nil {
				result.Hidden = ShouldHide(o, result, baseline)
			} else {
				result.Hidden = true // errors only surface under -v
			}
			results <- result
		}()
	}

	// Dispatch everything from a goroutine so the printer below can stream
	// results as they land instead of waiting for the whole fan-out.
	go func() {
		for _, url := range o.UrlList {
			baseline := computeBaseline(MyClient, url, o.PrimaryVerb, o.BasePath)
			if !baseline.OK {
				fmt.Printf("%s[!] No baseline for %s (control request failed) — showing all rows%s\n",
					WarningColor, url, EndColor)
			}
			pathURL := fmt.Sprintf("%s/%s", url, o.BasePath)

			if o.VerbBypasses {
				for _, method := range VerbBypasses {
					run(baseline, pathURL, method, "", "")
				}
			}

			if o.HeaderBypasses {
				for _, hb := range headerBypasses {
					target := url // rewrite headers hit root
					if hb.TargetPath {
						target = pathURL
					} else {
						target = fmt.Sprintf("%s/", url)
					}
					value := substitute(hb.Value, url, o.BasePath)
					run(baseline, target, o.PrimaryVerb, hb.Header, value)
				}
			}

			if o.PathBypasses {
				for _, pathFmtStr := range PathBypasses {
					finalURL := substitute(pathFmtStr, url, o.BasePath)
					run(baseline, finalURL, o.PrimaryVerb, "", "")
				}
			}

			// Path manipulation combined with header manipulation.
			if o.HeaderBypasses && o.PathBypasses {
				for _, pathFmtStr := range PathBypasses {
					finalURL := substitute(pathFmtStr, url, o.BasePath)
					for _, hb := range headerBypasses {
						if !hb.TargetPath {
							continue // rewrite headers are meaningless with a mangled path
						}
						value := substitute(hb.Value, url, o.BasePath)
						run(baseline, finalURL, o.PrimaryVerb, hb.Header, value)
					}
				}
			}
		}

		wg.Wait()
		close(results)
	}()

	// Stream results as they arrive.
	for result := range results {
		if result.Err != nil {
			if o.Verbose {
				fmt.Printf("%s[err]%s %-8s %-90s %s\n", ErrorColor, EndColor,
					result.Method, result.URL, result.Err)
			}
			continue
		}
		if result.Hidden && !(o.ShowAll || o.Verbose) {
			continue
		}
		fmt.Println(FormatResult(o, result))
	}

	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	timeout := atomic.LoadInt64(&MyClient.UserOptions.TimeoutRequests)
	succeeded := atomic.LoadInt64(&MyClient.UserOptions.TotalRequestsSucceeded)
	failed := atomic.LoadInt64(&MyClient.UserOptions.TotalRequestsFailed)
	fmt.Printf("\n\033[1;37m[Total Requests]\033[0m: %d\t\033[1;34m[Timed-out]\033[0m: %d\t\033[1;32m[Succeeded]\033[0m: %d\t\033[1;31m[Failed]\033[0m: %d\t\033[1;37m[Elapsed time]\033[0m: %.2f seconds\n",
		timeout+succeeded+failed,
		timeout, succeeded, failed,
		elapsedSeconds,
	)
}
