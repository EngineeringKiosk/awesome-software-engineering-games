// Package github contains helpers for working with the GitHub API in this
// project. The actual API calls are issued via github.com/google/go-github;
// this package only provides small utilities such as URL parsing.
package github

import (
	"net/url"
	"strings"
)

// ParseRepoURL extracts (owner, repo) from a GitHub HTTPS URL such as
// https://github.com/rlane/oort3 or https://github.com/rlane/oort3.git.
//
// It returns ok=false for any URL that is not on github.com or that does not
// look like an /<owner>/<repo> path. Callers should treat ok=false as
// "not a GitHub repository, skip silently".
func ParseRepoURL(raw string) (owner, repo string, ok bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", false
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", "", false
	}
	if !strings.EqualFold(u.Host, "github.com") {
		return "", "", false
	}

	path := strings.Trim(u.Path, "/")
	path = strings.TrimSuffix(path, ".git")

	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	return parts[0], parts[1], true
}
