package gobypasser

type Options struct {
	BaseURL    string
	FileOfUrls string
	BasePath   string
	UrlList    []string

	PathBypasses   bool
	HeaderBypasses bool
	VerbBypasses   bool

	UserAgent          string
	FilterResponseCode string
	FilterResponseSize string

	ParsedFilterResponseCode []string
	ParsedFilterResponseSize []string

	Colors  bool
	Verbose bool
	Timeout int

	TimeoutRequests        int
	TotalRequestsFailed    int
	TotalRequestsSucceeded int

	Threads int
}

func SetDefaultOptions(o *Options) {

	o.FileOfUrls = ""
	o.BaseURL = ""
	o.BasePath = ""

	o.Colors = true
	o.Verbose = false
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

}
