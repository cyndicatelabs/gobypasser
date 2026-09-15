package gobypasser

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

type HttpClient struct {
	UserOptions *Options
	HttpClient  *http.Client
}

// Result is the outcome of a single request. A nil Err means the request
// completed; StatusCode and Size are only meaningful in that case.
type Result struct {
	StatusCode int
	Size       int
	Method     string
	URL        string
	Header     string // custom bypass header shown in the table, or "N/A"
	Hidden     bool   // matches the baseline or a filter; suppressed unless -show-all/-v
	Err        error
}

func NewHttpClient(o *Options) HttpClient {
	var hc HttpClient

	// By default do not follow redirects: for bypass testing the raw 3xx and its
	// Location are the signal. -L (FollowRedirects) restores hop-following.
	checkRedirect := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	if o.FollowRedirects {
		checkRedirect = nil
	}

	hc.HttpClient = &http.Client{
		CheckRedirect: checkRedirect,
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

	req.Header.Set("User-Agent", MyClient.UserOptions.UserAgent)
	return req
}

func MakeHttpRequest(MyClient HttpClient, Request http.Request) *Result {

	result := &Result{
		Method: Request.Method,
		Header: HeaderToString(Request.Header),
	}
	if Request.URL != nil {
		result.URL = Request.URL.String()
	}

	res, err := MyClient.HttpClient.Do(&Request)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			atomic.AddInt64(&MyClient.UserOptions.TimeoutRequests, 1)
		} else {
			atomic.AddInt64(&MyClient.UserOptions.TotalRequestsFailed, 1)
		}
		result.Err = err
		return result
	}
	defer res.Body.Close()

	atomic.AddInt64(&MyClient.UserOptions.TotalRequestsSucceeded, 1)

	body, _ := io.ReadAll(res.Body)
	result.StatusCode = res.StatusCode
	result.Size = len(body)
	return result
}
