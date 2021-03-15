package gobypasser

import (
	"fmt"
	"strings"
	"sync"
)

func Start(o *Options) {

	fmt.Println()
	fmt.Print("Settings:")
	fmt.Println()
	fmt.Printf("  Base URL   : %s\n", o.BaseURL)
	fmt.Printf("  Base Path  : %s\n", o.BasePath)
	fmt.Printf("  User-Agent : %s\n", o.UserAgent)
	fmt.Println()

	var wg sync.WaitGroup
	var FinalURL string
	FinalURL = o.BaseURL + o.BasePath
	MyClient := NewHttpClient(o)

	PrintTableHeader()

	if o.VerbBypasses {
		for _, Method := range VerbBypasses {
			req := NewHttpRequest(MyClient, FinalURL, Method)
			go MakeHttpRequest(MyClient, req, &wg)
			wg.Add(1)
		}
		wg.Wait()
	}

	if o.HeaderBypasses {
		for _, Hdr := range HeaderBypassesHdr {
			for _, Val := range HeaderBypassesVal {
				req := NewHttpRequest(MyClient, FinalURL, "GET")

				Val = strings.ReplaceAll(Val, "{base_path}", MyClient.UserOptions.BasePath)
				Val = strings.ReplaceAll(Val, "{base_url}", MyClient.UserOptions.BaseURL)
				req.Header.Add(Hdr, Val)

				go MakeHttpRequest(MyClient, req, &wg)
				wg.Add(1)
			}
			wg.Wait()
		}
	}

	if o.PathBypasses {
		for _, PathFmtStr := range PathBypasses {
			PathFmtStr = strings.ReplaceAll(PathFmtStr, "{base_path}", MyClient.UserOptions.BasePath)
			PathFmtStr = strings.ReplaceAll(PathFmtStr, "{base_url}", MyClient.UserOptions.BaseURL)
			FinalURL = PathFmtStr

			req := NewHttpRequest(MyClient, FinalURL, "GET")
			go MakeHttpRequest(MyClient, req, &wg)
			wg.Add(1)
		}
		wg.Wait()
	}
}
