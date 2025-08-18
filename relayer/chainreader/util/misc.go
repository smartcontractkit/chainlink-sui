package util

import "unicode"

// SnakeToCamel converts "snake_case" to "SnakeCase".
func SnakeToCamel(s string) string {
	out := make([]rune, 0, len(s))
	upperNext := true
	for _, r := range s {
		if r == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			out = append(out, unicode.ToUpper(r))
			upperNext = false
		} else {
			out = append(out, unicode.ToLower(r))
		}
	}
	return string(out)
}

// CamelToSnake converts "camelCase" or "CamelCase" into "snake_case".
// It keeps acronyms together (e.g., "HTTPServer" -> "http_server")
// and inserts boundaries around digits (e.g., "version2Beta" -> "version2_beta").
func CamelToSnake(s string) string {
	r := []rune(s)
	n := len(r)
	if n == 0 {
		return ""
	}

	out := make([]rune, 0, n*2)

	for i := 0; i < n; i++ {
		cur := r[i]
		var prev, next rune
		if i > 0 {
			prev = r[i-1]
		}
		if i+1 < n {
			next = r[i+1]
		}

		// Normalize existing underscores to single underscores
		if cur == '_' {
			if len(out) > 0 && out[len(out)-1] != '_' {
				out = append(out, '_')
			}
			continue
		}

		// Determine whether to insert an underscore before cur.
		insertUnderscore := false

		if i > 0 {
			switch {
			// lower|digit → Upper   (e.g., aA, 3X)
			case unicode.IsUpper(cur) && (unicode.IsLower(prev) || unicode.IsDigit(prev)):
				insertUnderscore = true

			// Upper(acronym) → Upper followed by lower   (e.g., HTTPServer: before 'S')
			case unicode.IsUpper(cur) && unicode.IsUpper(prev) && i+1 < n && unicode.IsLower(next):
				insertUnderscore = true

			// letter → digit   (e.g., block2)
			case unicode.IsDigit(cur) && unicode.IsLetter(prev):
				insertUnderscore = true

			// digit → letter   (e.g., v2Beta: before 'B')
			case unicode.IsLetter(cur) && unicode.IsDigit(prev):
				insertUnderscore = true
			}
		}

		if insertUnderscore && (len(out) == 0 || out[len(out)-1] != '_') {
			out = append(out, '_')
		}

		out = append(out, unicode.ToLower(cur))
	}

	// Trim trailing underscore if produced by stray input underscores
	if len(out) > 0 && out[len(out)-1] == '_' {
		out = out[:len(out)-1]
	}
	return string(out)
}
