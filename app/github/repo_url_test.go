package github

import "testing"

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantOwner string
		wantRepo  string
		wantOK    bool
	}{
		{"basic https", "https://github.com/rlane/oort3", "rlane", "oort3", true},
		{"trailing slash", "https://github.com/rlane/oort3/", "rlane", "oort3", true},
		{".git suffix", "https://github.com/rlane/oort3.git", "rlane", "oort3", true},
		{"http scheme", "http://github.com/rlane/oort3", "rlane", "oort3", true},
		{"mixed case host", "https://GitHub.com/rlane/oort3", "rlane", "oort3", true},
		{"non-github host", "https://gitlab.com/rlane/oort3", "", "", false},
		{"only owner", "https://github.com/rlane", "", "", false},
		{"too many segments", "https://github.com/rlane/oort3/issues", "", "", false},
		{"empty", "", "", "", false},
		{"not a url", ":://broken", "", "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, ok := ParseRepoURL(tc.input)
			if owner != tc.wantOwner || repo != tc.wantRepo || ok != tc.wantOK {
				t.Errorf("ParseRepoURL(%q) = (%q, %q, %v), want (%q, %q, %v)",
					tc.input, owner, repo, ok, tc.wantOwner, tc.wantRepo, tc.wantOK)
			}
		})
	}
}
