package gobypasser

var HeaderBypassesHdr = []string{
	"Referer",
	"X-Forwarded-For",
	"X-Forwarded-Host",
	"X-Originating-IP",
	"X-rewrite-url",
	"X-Original-URL",
	"X-Remote-IP",
	"X-Client-IP",
	"X-Host",
}

var HeaderBypassesVal = []string{
	"{base_url}/{base_path}",
	"/{base_path}",
	"{base_path}",
	"http://127.0.0.1:80",
	"https://127.0.0.1:443",
	"http://127.0.0.1",
	"https://127.0.0.1",
	"127.0.0.1",
	"127.0.0.1:80",
	"127.0.0.1:433",
	"http://127.0.0.1:80/{base_path}",
	"https://127.0.0.1:443/{base_path}",
	"http://127.0.0.1/{base_path}",
	"https://127.0.0.1/{base_path}",
	"127.0.0.1/{base_path}",
	"127.0.0.1:80/{base_path}",
	"127.0.0.1:433/{base_path}",
	"[::1]/{base_path}",
	"[::1]:80/{base_path}",
	"[::1]:433/{base_path}",
	"localhost/{base_path}",
	"localhost:80/{base_path}",
	"localhost:433/{base_path}",
}
