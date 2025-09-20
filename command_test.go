package easycmd

import (
	"reflect"
	"testing"
)

func TestCommandName(t *testing.T) {
	tests := []struct {
		name     string
		command  command
		expected string
	}{
		{
			name:     "단순 명령어",
			command:  "ls",
			expected: "ls",
		},
		{
			name:     "인수가 있는 명령어",
			command:  "ls -la",
			expected: "ls",
		},
		{
			name:     "빈 명령어",
			command:  "",
			expected: "",
		},
		{
			name:     "공백만 있는 명령어",
			command:  "   ",
			expected: "",
		},
		{
			name:     "bash 래핑된 명령어",
			command:  "bash -c echo hello",
			expected: "bash",
		},
		{
			name:     "powershell 래핑된 명령어",
			command:  "powershell.exe Get-Process",
			expected: "powershell.exe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.Name()
			if result != tt.expected {
				t.Errorf("Name() = %q, 기대값: %q", result, tt.expected)
			}
		})
	}
}

func TestCommandArgs(t *testing.T) {
	tests := []struct {
		name     string
		command  command
		expected []string
	}{
		// 일반 명령어 테스트
		{
			name:     "인수 없는 명령어",
			command:  "ls",
			expected: []string{},
		},
		{
			name:     "단일 인수",
			command:  "ls -la",
			expected: []string{"-la"},
		},
		{
			name:     "여러 인수",
			command:  "git commit -m message",
			expected: []string{"commit", "-m", "message"},
		},
		{
			name:     "인용부호가 있는 인수",
			command:  "echo 'hello world'",
			expected: []string{"hello world"},
		},
		{
			name:     "빈 명령어",
			command:  "",
			expected: []string{},
		},

		// bash 래핑된 명령어 테스트
		{
			name:     "bash 단순 명령어",
			command:  "bash -c ls",
			expected: []string{"-c", "ls"},
		},
		{
			name:     "bash 복합 명령어",
			command:  "bash -c echo hello && ls",
			expected: []string{"-c", "echo hello && ls"},
		},
		{
			name:     "bash 파이프 명령어",
			command:  "bash -c cat file.txt | grep pattern",
			expected: []string{"-c", "cat file.txt | grep pattern"},
		},
		{
			name:     "bash 명령어에 bash 문자열이 포함된 경우",
			command:  "bash -c echo 'bash -c test'",
			expected: []string{"-c", "echo 'bash -c test'"},
		},

		// PowerShell 래핑된 명령어 테스트
		{
			name:     "powershell 단순 명령어",
			command:  "powershell.exe Get-Process",
			expected: []string{"-Command", "Get-Process"},
		},
		{
			name:     "powershell 복합 명령어",
			command:  "powershell.exe Get-Process | Where-Object {$_.Name -eq 'chrome'}",
			expected: []string{"-Command", "Get-Process | Where-Object {$_.Name -eq 'chrome'}"},
		},
		{
			name:     "powershell 명령어에 powershell 문자열이 포함된 경우",
			command:  "powershell.exe Write-Host 'powershell.exe test'",
			expected: []string{"-Command", "Write-Host 'powershell.exe test'"},
		},

		// 엣지 케이스
		{
			name:     "공백만 있는 명령어",
			command:  "   ",
			expected: []string{},
		},
		{
			name:     "bash 접두사만",
			command:  "bash -c ",
			expected: []string{"-c", ""},
		},
		{
			name:     "powershell 접두사만",
			command:  "powershell.exe ",
			expected: []string{"-Command", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.Args()

			// 빈 슬라이스와 nil 슬라이스를 동일하게 처리
			if len(result) == 0 && len(tt.expected) == 0 {
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Args() = %v, 기대값: %v", result, tt.expected)
			}
		})
	}
}

func TestCommandShellCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  command
		expected command
	}{
		{
			name:     "단순 명령어",
			command:  "ls",
			expected: "bash -c ls",
		},
		{
			name:     "복합 명령어",
			command:  "echo hello && ls",
			expected: "bash -c echo hello && ls",
		},
		{
			name:     "빈 명령어",
			command:  "",
			expected: "bash -c ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.ShellCommand()
			if result != tt.expected {
				t.Errorf("ShellCommand() = %q, 기대값: %q", result, tt.expected)
			}
		})
	}
}

func TestCommandPowershellCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  command
		expected command
	}{
		{
			name:     "단순 명령어",
			command:  "Get-Process",
			expected: "powershell.exe Get-Process",
		},
		{
			name:     "복합 명령어",
			command:  "Get-Process | Where-Object {$_.Name -eq 'chrome'}",
			expected: "powershell.exe Get-Process | Where-Object {$_.Name -eq 'chrome'}",
		},
		{
			name:     "빈 명령어",
			command:  "",
			expected: "powershell.exe ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.PowershellCommand()
			if result != tt.expected {
				t.Errorf("PowershellCommand() = %q, 기대값: %q", result, tt.expected)
			}
		})
	}
}

func TestCommandString(t *testing.T) {
	tests := []struct {
		name     string
		command  command
		expected string
	}{
		{
			name:     "단순 명령어",
			command:  "ls",
			expected: "ls",
		},
		{
			name:     "복합 명령어",
			command:  "echo hello && ls",
			expected: "echo hello && ls",
		},
		{
			name:     "빈 명령어",
			command:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.command.String()
			if result != tt.expected {
				t.Errorf("String() = %q, 기대값: %q", result, tt.expected)
			}
		})
	}
}
