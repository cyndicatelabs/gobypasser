package gobypasser

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	SuccessColor = "\033[1;32m"
	InfoColor    = "\033[1;34m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	EndColor     = "\033[0m"
)

func PrintTableHeader() {
	fmt.Printf("%-15s %-15s %-20s %-90s %-60s\n", "Response Code", "Response Size", "Verb", "Path", "Custom Header")
	fmt.Printf("%s\n", strings.Repeat("_", 170))

}

func HeaderToString(Headers http.Header) string {

	var HdrString = "N/A"

	for k, v := range Headers {
		if stringInSlice(k, HeaderBypassesHdr) {
			HdrString = fmt.Sprintf("%s: %s", k, v[0])
		}
	}
	return HdrString
}

func PrintResult(MyClient HttpClient, Request http.Request, Response http.Response) {

	defer Response.Body.Close()
	var Color string = EndColor
	var strCode = strconv.Itoa(Response.StatusCode)

	if strCode[0] == '2' {
		Color = SuccessColor
	} else if strCode[0] == '3' {
		Color = InfoColor
	} else if strCode[0] == '4' {
		Color = ErrorColor
	} else if strCode[0] == '5' {
		Color = WarningColor
	} else {
		Color = EndColor
	}

	var strLength = "0"
	body, err := io.ReadAll(Response.Body)
	if err != nil {
		strLength = "-1"
	}
	strLength = strconv.Itoa(len(body))

	if !(stringInSlice(strCode, MyClient.UserOptions.ParsedFilterResponseCode) ||
		stringInSlice(strLength, MyClient.UserOptions.ParsedFilterResponseSize)) ||
		(len(MyClient.UserOptions.ParsedFilterResponseSize) == 1 &&
			len(MyClient.UserOptions.ParsedFilterResponseCode) == 1) {

		fmt.Printf(
			"%s%-15d%s %-15s %-20s %-90s %-60v\n",
			Color,
			Response.StatusCode,
			EndColor,
			strLength,
			Request.Method,
			Request.URL,
			HeaderToString(Request.Header),
		)
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.EqualFold(b, a) {
			return true
		}
	}
	return false
}
