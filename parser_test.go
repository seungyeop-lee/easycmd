package easycmd

import (
	"reflect"
	"testing"
)

func TestParseCommandArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		// 기본 케이스
		{
			name:     "단일 명령어",
			input:    "ls",
			expected: []string{"ls"},
		},
		{
			name:     "여러 인수",
			input:    "ls -la /home",
			expected: []string{"ls", "-la", "/home"},
		},
		{
			name:     "빈 문자열",
			input:    "",
			expected: []string{},
		},

		// 공백 처리
		{
			name:     "앞뒤 공백",
			input:    "  ls  ",
			expected: []string{"ls"},
		},
		{
			name:     "여러 공백",
			input:    "ls   -la",
			expected: []string{"ls", "-la"},
		},
		{
			name:     "공백만",
			input:    "   ",
			expected: []string{},
		},

		// 단일 인용부호 처리
		{
			name:     "단일 인용부호로 감싼 문자열",
			input:    "echo 'hello world'",
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "빈 단일 인용부호",
			input:    "echo ''",
			expected: []string{"echo"},
		},
		{
			name:     "단일 인용부호 안의 공백",
			input:    "echo 'hello   world'",
			expected: []string{"echo", "hello   world"},
		},

		// 이중 인용부호 처리
		{
			name:     "이중 인용부호로 감싼 문자열",
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
		},
		{
			name:     "빈 이중 인용부호",
			input:    `echo ""`,
			expected: []string{"echo"},
		},
		{
			name:     "이중 인용부호 안의 공백",
			input:    `echo "hello   world"`,
			expected: []string{"echo", "hello   world"},
		},

		// 혼합 인용부호
		{
			name:     "단일 인용부호 안의 이중 인용부호",
			input:    `echo 'hello "quoted" world'`,
			expected: []string{"echo", `hello "quoted" world`},
		},
		{
			name:     "이중 인용부호 안의 단일 인용부호",
			input:    `echo "hello 'quoted' world"`,
			expected: []string{"echo", "hello 'quoted' world"},
		},

		// 엣지 케이스
		{
			name:     "인용부호만",
			input:    `''`,
			expected: []string{},
		},
		{
			name:     "이중 인용부호만",
			input:    `""`,
			expected: []string{},
		},
		{
			name:     "닫히지 않은 단일 인용부호",
			input:    "echo 'hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "닫히지 않은 이중 인용부호",
			input:    `echo "hello`,
			expected: []string{"echo", "hello"},
		},
		{
			name:     "연속된 인용부호",
			input:    "echo ''world''",
			expected: []string{"echo", "world"},
		},

		// 실제 사용 사례
		{
			name:     "Git 커밋 메시지",
			input:    "git commit -m 'Initial commit'",
			expected: []string{"git", "commit", "-m", "Initial commit"},
		},
		{
			name:     "공백이 있는 파일 경로",
			input:    "cat '/path/with spaces/file.txt'",
			expected: []string{"cat", "/path/with spaces/file.txt"},
		},
		{
			name:     "복잡한 명령어",
			input:    `find . -name "*.go" -exec grep "func" {} \;`,
			expected: []string{"find", ".", "-name", "*.go", "-exec", "grep", "func", "{}", `\;`},
		},
		{
			name:     "여러 옵션과 인용부호",
			input:    `docker run -it --name "my container" ubuntu:latest`,
			expected: []string{"docker", "run", "-it", "--name", "my container", "ubuntu:latest"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommandArgs(tt.input)

			// 빈 슬라이스와 nil 슬라이스를 동일하게 처리
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseCommandArgs(%q) = %v, 기대값: %v", tt.input, result, tt.expected)
			}
		})
	}
}
