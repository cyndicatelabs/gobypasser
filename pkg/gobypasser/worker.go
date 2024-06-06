package gobypasser

import (
	"fmt"
	"strings"
	"sync"
)

func Start(o *Options) {

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

	for _, url := range o.UrlList {
		FinalURL = url + o.BasePath

		if o.VerbBypasses {
			for _, Method := range VerbBypasses {
				wg.Add(1)
				go func(url, method string) {
					defer wg.Done()
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
					go func(url, hdr, val string) {
						defer wg.Done()
						req := NewHttpRequest(MyClient, url, "GET")
						val = strings.ReplaceAll(val, "{base_path}", MyClient.UserOptions.BasePath)
						val = strings.ReplaceAll(val, "{base_url}", MyClient.UserOptions.BaseURL)
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
				go func(url, pathFmtStr string) {
					defer wg.Done()                                                // Create a copy of the loop variable
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
	}

	// Close the results channel when all requests are done
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}

}
