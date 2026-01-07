package github

import "testing"

func TestLimitDiffFiles(t *testing.T) {
	diff := "" +
		"diff --git a/a.txt b/a.txt\n" +
		"index 1..2 100644\n" +
		"--- a/a.txt\n" +
		"+++ b/a.txt\n" +
		"@@ -1 +1 @@\n" +
		"-a\n" +
		"+aa\n" +
		"diff --git a/b.txt b/b.txt\n" +
		"--- a/b.txt\n" +
		"+++ b/b.txt\n" +
		"diff --git a/c.txt b/c.txt\n" +
		"--- a/c.txt\n" +
		"+++ b/c.txt\n"

	out, truncated := LimitDiffFiles(diff, 2)
	if !truncated {
		t.Fatalf("expected truncated=true")
	}
	if count := countPrefixLines(out, "diff --git "); count != 2 {
		t.Fatalf("expected 2 diff sections, got %d", count)
	}
}

func countPrefixLines(s, prefix string) int {
	n := 0
	start := 0
	for start < len(s) {
		end := start
		for end < len(s) && s[end] != '\n' {
			end++
		}
		line := s[start:end]
		if len(line) >= len(prefix) && line[:len(prefix)] == prefix {
			n++
		}
		if end == len(s) {
			break
		}
		start = end + 1
	}
	return n
}
