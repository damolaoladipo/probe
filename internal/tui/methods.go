package tui

var methods = []string{
	"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
}

// nextMethod returns the next HTTP method in GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS.
// Unknown values wrap to GET.
func nextMethod(current string) string {
	for i, m := range methods {
		if m == current {
			return methods[(i+1)%len(methods)]
		}
	}
	return "GET"
}
