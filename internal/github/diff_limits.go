package github

import (
	"bufio"
	"strings"
)

func LimitDiffFiles(diff string, maxFiles int) (string, bool) {
	if maxFiles <= 0 {
		return diff, false
	}

	sc := bufio.NewScanner(strings.NewReader(diff))

	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 4*1024*1024)

	var b strings.Builder
	var files int

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, "diff --git ") {
			files++
			if files > maxFiles {
				return b.String(), true
			}
		}

		b.WriteString(line)
		b.WriteByte('\n')
	}

	if err := sc.Err(); err != nil {
		return b.String(), true
	}

	return b.String(), false
}
