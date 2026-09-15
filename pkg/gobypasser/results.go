package gobypasser

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	SuccessColor = "\033[1;32m"
	InfoColor    = "\033[1;34m"
	WarningColor = "\033[1;33m"
	ErrorColor   = "\033[1;31m"
	DimColor     = "\033[2m"
	EndColor     = "\033[0m"
)

// Baseline is the response to the unmodified request for a target. A bypass
// whose response matches OK && StatusCode && Size is considered noise.
type Baseline struct {
	StatusCode int
	Size       int
	OK         bool
}

func (b Baseline) Matches(r *Result) bool {
	return b.OK && r.StatusCode == b.StatusCode && r.Size == b.Size
}

// colorEnabled reports whether ANSI colors should be emitted: honor -no-color
// and the NO_COLOR convention, and only colorize when stdout is a terminal.
func colorEnabled(o *Options) bool {
	if o.NoColor || os.Getenv("NO_COLOR") != "" {
		return false
	}
	info, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func PrintTableHeader() {
	fmt.Printf("%-15s %-15s %-20s %-90s %-60s\n", "Response Code", "Response Size", "Verb", "Path", "Custom Header")
	fmt.Printf("%s\n", strings.Repeat("_", 170))
}

// HeaderToString renders the single bypass header set on a request, or "N/A".
func HeaderToString(Headers http.Header) string {
	names := AllBypassHeaderNames()
	for _, name := range names {
		if v := Headers.Get(name); v != "" {
			return fmt.Sprintf("%s: %s", http.CanonicalHeaderKey(name), v)
		}
	}
	return "N/A"
}

func statusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return SuccessColor
	case code >= 300 && code < 400:
		return InfoColor
	case code >= 400 && code < 500:
		return ErrorColor
	case code >= 500:
		return WarningColor
	default:
		return EndColor
	}
}

// ShouldHide decides whether a completed result is noise: it matches the
// baseline, or its code/size is on a manual -hc/-hs hide list.
func ShouldHide(o *Options, r *Result, baseline Baseline) bool {
	if baseline.Matches(r) {
		return true
	}
	if StrInSlice(strconv.Itoa(r.StatusCode), o.ParsedFilterResponseCode) {
		return true
	}
	if StrInSlice(strconv.Itoa(r.Size), o.ParsedFilterResponseSize) {
		return true
	}
	return false
}

// FormatResult renders one table row. Hidden rows (shown only under
// -show-all/-v) are dimmed so real hits still stand out.
func FormatResult(o *Options, r *Result) string {
	useColor := colorEnabled(o)

	code := statusColor(r.StatusCode)
	end := EndColor
	if !useColor {
		code, end = "", ""
	}

	line := fmt.Sprintf(
		"%s%-15d%s %-15d %-20s %-90s %-60s",
		code, r.StatusCode, end,
		r.Size,
		r.Method,
		r.URL,
		r.Header,
	)

	if r.Hidden && useColor {
		line = DimColor + line + EndColor
	}
	return line
}
