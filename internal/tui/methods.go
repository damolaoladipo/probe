package tui

var methods = []string{
	"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
}

func nextMethod(current string) string {
	for i, m := range methods {
		if m == current {
			return methods[(i+1)%len(methods)]
		}
	}
	return "GET"
}
