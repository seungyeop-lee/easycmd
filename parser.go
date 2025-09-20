package easycmd

import "strings"

// parseCommandArgs 명령어 문자열을 인수 배열로 파싱 (인용부호 처리 포함)
func parseCommandArgs(cmd string) []string {
	var args []string
	var currentToken strings.Builder
	var insideQuotes bool
	var activeQuoteChar rune

	for _, char := range cmd {
		if !insideQuotes && isQuoteChar(char) {
			// 인용부호 시작
			insideQuotes = true
			activeQuoteChar = char
		} else if insideQuotes && char == activeQuoteChar {
			// 인용부호 종료
			insideQuotes = false
		} else if !insideQuotes && char == ' ' {
			// 공백으로 토큰 분리
			args = addTokenIfNotEmpty(args, currentToken.String())
			currentToken.Reset()
		} else {
			// 일반 문자 추가
			currentToken.WriteRune(char)
		}
	}

	// 마지막 토큰 추가
	args = addTokenIfNotEmpty(args, currentToken.String())
	return args
}

// isQuoteChar 인용부호 문자인지 확인
// 예: isQuoteChar('"') -> true, isQuoteChar('a') -> false
func isQuoteChar(r rune) bool {
	return r == '"' || r == '\''
}

// addTokenIfNotEmpty 토큰이 비어있지 않으면 배열에 추가
// 예: addTokenIfNotEmpty([]string{"ls"}, "file") -> []string{"ls", "file"}
func addTokenIfNotEmpty(args []string, token string) []string {
	if token != "" {
		return append(args, token)
	}
	return args
}
