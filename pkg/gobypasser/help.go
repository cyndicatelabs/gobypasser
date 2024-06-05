package gobypasser

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
)

type UsageSection struct {
	Name          string
	Description   string
	Flags         []UsageFlag
	ExpectedFlags []string
}

type UsageFlag struct {
	Name        string
	Description string
	Default     string
	Required    bool
}

func Banner() {
	Version := "0.1"
	fmt.Println()
	fmt.Printf("\t\tGoBypasser %s - https://www.github.com/cyndicatelabs/gobypasser -\n", Version)
	fmt.Println()
	fmt.Printf("\t    A tool to aid finding URL bypasses using a number of different techniques.\n")
	fmt.Println()
}

func (u *UsageSection) PrintSection(max_length int) {
	fmt.Printf("%s:\n", u.Name)
	for _, f := range u.Flags {
		f.PrintFlag(max_length)
	}
	fmt.Printf("\n")
}

func (f *UsageFlag) PrintFlag(max_length int) {
	format := fmt.Sprintf("    -%%-%ds %%s", max_length)
	if f.Default != "" {
		format = format + " (default: %s)\n"
		fmt.Printf(format, f.Name, f.Description, f.Default)
	} else {
		format = format + "\n"
		fmt.Printf(format, f.Name, f.Description)
	}
}

func Usage() {
	Banner()

	uGeneral := UsageSection{
		Name:          "General Options",
		Description:   "",
		Flags:         make([]UsageFlag, 0),
		ExpectedFlags: []string{"u", "p", "v", "t"},
	}
	uAttack := UsageSection{
		Name:          "Attack Options",
		Description:   "Bypass method flags that will be used on execution.",
		Flags:         make([]UsageFlag, 0),
		ExpectedFlags: []string{"verbs", "headers", "paths"},
	}
	uFilter := UsageSection{
		Name:          "Filter Options",
		Description:   "Filters for the response filtering.",
		Flags:         make([]UsageFlag, 0),
		ExpectedFlags: []string{"hs", "hc"},
	}

	sections := []UsageSection{uGeneral, uAttack, uFilter}

	max_length := 0
	flag.VisitAll(func(f *flag.Flag) {
		for i, section := range sections {
			if strInSlice(f.Name, section.ExpectedFlags) {
				sections[i].Flags = append(sections[i].Flags, UsageFlag{
					Name:        f.Name,
					Description: f.Usage,
					Default:     f.DefValue,
				})
			}
		}
		if len(f.Name) > max_length {
			max_length = len(f.Name)
		}
	})

	for _, section := range sections {
		section.PrintSection(max_length)
	}

	fmt.Printf("EXAMPLE USAGE:\n")

	fmt.Println("\tgobypasser -u https://www.google.com/ -p /api")
	fmt.Println("\tgobypasser -u https://www.google.com/ -p /api --hc 302,404 --verbs")
	fmt.Println("\tgobypasser -u https://www.google.com/ -p /api --hs 1024 --headers")
	fmt.Println()
}

func strInSlice(val string, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

func VerifyOptions(o *Options) (bool, string) {

	_, err := url.ParseRequestURI(o.BaseURL)
	if err != nil {
		return false, fmt.Sprintf("Missing URL parameter (-u): %s", err)
	}

	if len(o.BasePath) == 0 {
		return false, fmt.Sprintf("Missing path parameter (-p)")
	}

	if !(o.HeaderBypasses || o.PathBypasses || o.VerbBypasses) {
		return false, fmt.Sprintf("Missing attack options")
	}

	o.BaseURL = strings.TrimSuffix(o.BaseURL, "/")

	if !(strings.HasPrefix(o.BasePath, "/")) {
		o.BasePath = "/" + o.BasePath
	}

	if len(o.FilterResponseCode) > 0 {
		for _, c := range strings.Split(o.FilterResponseCode, ",") {
			o.ParsedFilterResponseCode = append(o.ParsedFilterResponseCode, c)
		}
	}

	if len(o.FilterResponseSize) > 0 {
		for _, c := range strings.Split(o.FilterResponseSize, ",") {
			o.ParsedFilterResponseSize = append(o.ParsedFilterResponseSize, c)
		}
	}

	return true, ""
}
