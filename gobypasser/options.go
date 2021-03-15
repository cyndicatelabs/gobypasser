package gobypasser

type Options struct {
	BaseURL  string
	BasePath string

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
}

func SetDefaultOptions(o *Options) {

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
}
