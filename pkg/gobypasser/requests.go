package gobypasser

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient struct {
	UserOptions *Options
	HttpClient  *http.Client
}

func NewHttpClient(o *Options) HttpClient {
	var hc HttpClient

	hc.HttpClient = &http.Client{
		CheckRedirect: nil,
		Timeout:       time.Duration(time.Duration(o.Timeout) * time.Second),
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 500,
			MaxConnsPerHost:     500,
			DialContext: (&net.Dialer{
				Timeout: time.Duration(time.Duration(o.Timeout) * time.Second),
			}).DialContext,
			TLSHandshakeTimeout: time.Duration(time.Duration(o.Timeout) * time.Second),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				Renegotiation:      tls.RenegotiateOnceAsClient,
			},
		},
	}

	hc.UserOptions = o
	return hc
}

func NewHttpRequest(MyClient HttpClient, FinalUrl string, Method string) http.Request {
	reqURL, _ := url.Parse(FinalUrl)
	var req http.Request

	if Method == "POST" {
		req = http.Request{
			Method: Method,
			URL:    reqURL,
			Header: http.Header{},
			Body:   io.NopCloser(strings.NewReader(`{"id":"1"}`)),
		}
	} else {
		req = http.Request{
			Method: Method,
			URL:    reqURL,
			Header: http.Header{},
		}
	}

	req.Header.Add("User-Agent", MyClient.UserOptions.UserAgent)
	return req
}

func MakeHttpRequest(MyClient HttpClient, Request http.Request) string {

	// defer wg.Done()
	Request.Header.Add("User-Agent", MyClient.UserOptions.UserAgent)

	res, err := MyClient.HttpClient.Do(&Request)
	if err != nil {
		// If its a timeout error
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			// fmt.Printf("Host timeout: %s\n - skipping...", Request.URL)
			MyClient.UserOptions.TimeoutRequests++
		} else {
			MyClient.UserOptions.TotalRequestsFailed++
		}
	} else {
		MyClient.UserOptions.TotalRequestsSucceeded++
		defer res.Body.Close()
	}

	if err == nil {
		return GetResult(MyClient, Request, *res)
	} else {
		return ""
	}
}
