package github

func TruncateUTF8ByBytes(s string, maxBytes int) (out string, truncated bool) {
	if maxBytes <= 0 {
		return "", len(s) > 0
	}

	if len(s) <= maxBytes {
		return s, false
	}

	var cutPos int
	for i := range s {
		if i > maxBytes {
			break
		}
		cutPos = i
	}

	if cutPos == 0 {
		return "", true
	}

	return s[:cutPos], true
}
