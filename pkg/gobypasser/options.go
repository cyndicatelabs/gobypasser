package gobypasser

type Options struct {
	BaseURL    string
	FileOfUrls string
	BasePath   string
	UrlList    []string

	PathBypasses   bool
	HeaderBypasses bool
	VerbBypasses   bool

	PrimaryVerb        string
	UserAgent          string
	FilterResponseCode string
	FilterResponseSize string

	ParsedFilterResponseCode []string
	ParsedFilterResponseSize []string

	NoColor         bool
	ShowAll         bool
	Verbose         bool
	FollowRedirects bool
	Timeout         int

	// Request counters. Incremented from many goroutines, so accessed
	// exclusively through sync/atomic.
	TimeoutRequests        int64
	TotalRequestsFailed    int64
	TotalRequestsSucceeded int64

	Threads int
}

func SetDefaultOptions(o *Options) {

	o.FileOfUrls = ""
	o.BaseURL = ""
	o.BasePath = ""

	o.PrimaryVerb = "GET"
	o.NoColor = false
	o.ShowAll = false
	o.Verbose = false
	o.FollowRedirects = false
	o.Timeout = 30

	o.PathBypasses = false
	o.HeaderBypasses = false
	o.VerbBypasses = false

	o.FilterResponseCode = ""
	o.FilterResponseSize = ""

	o.ParsedFilterResponseCode = []string{}
	o.ParsedFilterResponseSize = []string{}

	o.TimeoutRequests = 0
	o.TotalRequestsSucceeded = 0
	o.TotalRequestsFailed = 0

	o.Threads = 30

}
