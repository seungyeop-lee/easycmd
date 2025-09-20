package easycmd

// parseCommandArgs parses command string into arguments, handling quoted strings
func parseCommandArgs(cmd string) []string {
	var args []string
	var current string
	var inQuotes bool
	var quoteChar rune

	for _, r := range cmd {
		switch {
		case !inQuotes && (r == '"' || r == '\''):
			inQuotes = true
			quoteChar = r
		case inQuotes && r == quoteChar:
			inQuotes = false
		case !inQuotes && r == ' ':
			if current != "" {
				args = append(args, current)
				current = ""
			}
		default:
			current += string(r)
		}
	}

	if current != "" {
		args = append(args, current)
	}

	return args
}
