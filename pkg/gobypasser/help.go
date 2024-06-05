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
	Version := "1.0"
	fmt.Println()
	fmt.Printf("\t\t\033[1;32mGoBypasser %s - https://www.github.com/cyndicatelabs/gobypasser - @cyndicatelabs\033[0m\n", Version)
	fmt.Println()
	fmt.Printf("A tool to help find 403 URL bypasses using a number of different techniques.\n")
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
		ExpectedFlags: []string{"u", "p", "v", "t", "f"},
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
			if StrInSlice(f.Name, section.ExpectedFlags) {
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

	fmt.Println("\tgobypasser -u https://www.google.com/ -p admin -verbs -hc 403")
	fmt.Println("\tgobypasser -u https://www.google.com/ -p admin -verbs -headers -paths -hc 403,400")
	fmt.Println("\tgobypasser -f ./url_list.txt -p admin -verbs -headers -paths -hc 403,400")
	fmt.Println()
}

func VerifyOptions(o *Options) (bool, string) {

	if len(o.FileOfUrls) > 0 {

		if _, err := FileExists(o.FileOfUrls); err != nil {
			return false, fmt.Sprintf("File does not exist: %s", o.BaseURL)
		}
		// parse the file and put the urls into the UrlList if they pass the ParseRequestURI check
		tmpUrls, err := ParseFile(o.FileOfUrls)
		if err != nil {
			return false, fmt.Sprintf("Error parsing file: %s", err)
		}

		for _, urlItem := range tmpUrls {
			if _, err := url.ParseRequestURI(urlItem); err == nil {
				urlItem = strings.TrimSuffix(urlItem, "/")
				o.UrlList = append(o.UrlList, urlItem)
			}
		}

	} else {
		_, err := url.ParseRequestURI(o.BaseURL)
		if err != nil {
			return false, fmt.Sprintf("Missing URL parameter (-u): %s", err)
		}
		o.BaseURL = strings.TrimSuffix(o.BaseURL, "/")
		o.UrlList = append(o.UrlList, o.BaseURL)
	}

	fmt.Printf("%+v\n", o.UrlList)

	if len(o.BasePath) == 0 {
		return false, fmt.Sprintln("Missing path parameter (-p)")
	}

	if !(strings.HasPrefix(o.BasePath, "/")) {
		o.BasePath = "/" + o.BasePath
	}

	if !(o.HeaderBypasses || o.PathBypasses || o.VerbBypasses) {
		return false, fmt.Sprintln("Missing attack options")
	}

	if len(o.FilterResponseCode) > 0 {
		o.ParsedFilterResponseCode = append(o.ParsedFilterResponseCode, strings.Split(o.FilterResponseCode, ",")...)
	}

	if len(o.FilterResponseSize) > 0 {
		o.ParsedFilterResponseSize = append(o.ParsedFilterResponseSize, strings.Split(o.FilterResponseSize, ",")...)
	}

	return true, ""
}
