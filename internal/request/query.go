package request

import (
	"fmt"
	"net/url"
	"strings"
)

// ApplyQuery merges query textarea lines onto rawURL and returns the new URL string.
// Each non-empty line must be key=value. Empty Query leaves the URL (and any ?query) unchanged.
// Existing keys in the URL are overwritten by q.Set, not appended.
func ApplyQuery(rawURL, queryLines string) (string, error) {

	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", err
	}

	q := u.Query()
	for _, line := range strings.Split(queryLines, "\n") {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}
		k, v, ok := strings.Cut(line, "=")

		if !ok {
			return "", fmt.Errorf("query line %q: want key=value", line)
		}
		q.Set(strings.TrimSpace(k), strings.TrimSpace(v))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
