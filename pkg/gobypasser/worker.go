package gobypasser

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func Start(o *Options) {

	start := time.Now()

	fmt.Println("Settings:")
	if o.BaseURL != "" {
		fmt.Printf("  Base URL   : %s\n", o.BaseURL)
	} else if o.FileOfUrls != "" {
		fmt.Printf("  Number of URLs : %d\n", len(o.UrlList))
	}
	fmt.Printf("  Base Path  : %s\n", o.BasePath)
	fmt.Printf("  User-Agent : %s\n", o.UserAgent)
	fmt.Println()

	var wg sync.WaitGroup
	var FinalURL string
	MyClient := NewHttpClient(o)

	PrintTableHeader()

	results := make(chan string, len(o.UrlList)*(len(VerbBypasses)+len(HeaderBypassesHdr)*len(HeaderBypassesVal)+len(PathBypasses)))
	workerPool := make(chan struct{}, o.Threads)

	for _, url := range o.UrlList {
		// FinalURL Should be url + o.BasePath + "/"
		FinalURL = fmt.Sprintf("%s/%s", url, o.BasePath)

		if o.VerbBypasses {
			for _, Method := range VerbBypasses {
				wg.Add(1)
				workerPool <- struct{}{}
				go func(url, method string) {
					defer wg.Done()
					defer func() { <-workerPool }()
					req := NewHttpRequest(MyClient, url, method)
					result := MakeHttpRequest(MyClient, req)
					if result != "" {
						results <- result
					}
				}(FinalURL, Method)
			}
		}

		if o.HeaderBypasses {
			for _, Hdr := range HeaderBypassesHdr {
				for _, Val := range HeaderBypassesVal {
					wg.Add(1)
					workerPool <- struct{}{}
					go func(url, hdr, val string) {
						defer wg.Done()
						defer func() { <-workerPool }()
						req := NewHttpRequest(MyClient, url, "GET")
						val = strings.ReplaceAll(val, "{base_path}", MyClient.UserOptions.BasePath)
						val = strings.ReplaceAll(val, "{base_url}", fmt.Sprintf(MyClient.UserOptions.BaseURL, "/"))
						req.Header.Add(hdr, val)
						result := MakeHttpRequest(MyClient, req)
						if result != "" {
							results <- result
						}
					}(FinalURL, Hdr, Val)
				}
			}
		}

		if o.PathBypasses {
			for _, PathFmtStr := range PathBypasses {
				wg.Add(1)
				workerPool <- struct{}{}
				go func(url, pathFmtStr string) {
					defer wg.Done() // Create a copy of the loop variable
					defer func() { <-workerPool }()
					pathFmtStr = strings.ReplaceAll(pathFmtStr, "{base_url}", url) // Use the copy instead of the loop variable
					pathFmtStr = strings.ReplaceAll(pathFmtStr, "{base_path}", MyClient.UserOptions.BasePath)
					finalURL := pathFmtStr
					req := NewHttpRequest(MyClient, finalURL, "GET")
					result := MakeHttpRequest(MyClient, req)
					if result != "" {
						results <- result
					}
				}(url, PathFmtStr)
			}
		}

		// Do both path manipulation and header manipulation together
		if o.HeaderBypasses && o.PathBypasses {
			for _, PathFmtStr := range PathBypasses {
				for _, Hdr := range HeaderBypassesHdr {
					for _, Val := range HeaderBypassesVal {
						wg.Add(1)
						workerPool <- struct{}{}
						go func(url, pathFmtStr, hdr, val string) {
							defer wg.Done()
							defer func() { <-workerPool }()
							pathFmtStr = strings.ReplaceAll(pathFmtStr, "{base_url}", url)
							pathFmtStr = strings.ReplaceAll(pathFmtStr, "{base_path}", MyClient.UserOptions.BasePath)
							finalURL := pathFmtStr
							req := NewHttpRequest(MyClient, finalURL, "GET")
							val = strings.ReplaceAll(val, "{base_path}", MyClient.UserOptions.BasePath)
							val = strings.ReplaceAll(val, "{base_url}", fmt.Sprintf(MyClient.UserOptions.BaseURL, "/"))
							req.Header.Add(hdr, val)
							result := MakeHttpRequest(MyClient, req)
							if result != "" {
								results <- result
							}
						}(url, PathFmtStr, Hdr, Val)
					}
				}
			}
		}
	}

	// Close the results channel when all requests are done
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

	// Print all timeout requests
	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	fmt.Printf("\n\033[1;37m[Total Requests]\033[0m: %d\t\033[1;34m[Timed-out]\033[0m: %d\t\033[1;32m[Succeeded]\033[0m: %d\t\033[1;31m[Failed]\033[0m: %d\t\033[1;37m[Elapsed time]\033[0m: %.2f seconds\n",
		MyClient.UserOptions.TimeoutRequests+MyClient.UserOptions.TotalRequestsSucceeded,
		MyClient.UserOptions.TimeoutRequests, MyClient.UserOptions.TotalRequestsSucceeded,
		MyClient.UserOptions.TotalRequestsFailed,
		elapsedSeconds,
	)
}
