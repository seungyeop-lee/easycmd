package test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/seungyeop-lee/easycmd"
)

func TestSimple(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when
	err := cmd.Run("echo hello world")

	// then
	if out.String() != "hello world\n" {
		t.Errorf("expected hello world, got %s", out.String())
	}
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRunShell(t *testing.T) {
	cmd := easycmd.New()

	err := cmd.RunShell("(cd .. && pwd)")

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRunMultiLineShell(t *testing.T) {
	cmd := easycmd.New()

	err := cmd.RunShell(`
	pwd
	ls -al
`)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRunWithDir(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when
	err := cmd.RunWithDir("pwd", "..")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	// 상위 디렉토리의 pwd 결과를 확인
	result := out.String()
	if !strings.Contains(result, "seungyeop-lee") {
		t.Errorf("expected path to contain 'seungyeop-lee', got %s", result)
	}
}

func TestRunShellWithDir(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - 상위 디렉토리에서 현재 디렉토리명 확인
	err := cmd.RunShellWithDir("basename $(pwd)", "..")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "easycmd" {
		t.Errorf("expected 'easycmd', got '%s'", result)
	}
}

func TestEmptyCommand(t *testing.T) {
	// given
	cmd := easycmd.New()

	// when
	err := cmd.Run("")

	// then
	if err == nil {
		t.Error("expected error for empty command, got nil")
	}

	if !errors.Is(err, easycmd.EmptyCmdError) {
		t.Errorf("expected EmptyCmdError, got %v", err)
	}
}

func TestInvalidCommand(t *testing.T) {
	// given
	cmd := easycmd.New()

	// when
	err := cmd.Run("nonexistentcommand12345")

	// then
	if err == nil {
		t.Error("expected error for invalid command, got nil")
		return
	}

	// 에러 메시지에 "명령어를 시작할 수 없습니다" 또는 "명령어 실행이 실패했거나"가 포함되어야 함
	errMsg := err.Error()
	if !strings.Contains(errMsg, "명령어를 시작할 수 없습니다") && !strings.Contains(errMsg, "명령어 실행이 실패했거나") {
		t.Errorf("expected error message to contain command execution error, got: %s", errMsg)
	}
}

func TestInvalidDirectory(t *testing.T) {
	// given
	cmd := easycmd.New()

	// when
	err := cmd.RunWithDir("echo test", "/nonexistent/directory/path12345")

	// then
	if err == nil {
		t.Error("expected error for invalid directory, got nil")
		return
	}

	// 에러 메시지 확인
	errMsg := err.Error()
	if !strings.Contains(errMsg, "명령어를 시작할 수 없습니다") && !strings.Contains(errMsg, "명령어 실행이 실패했거나") {
		t.Errorf("expected error message to contain command execution error, got: %s", errMsg)
	}
}

func TestStdErr(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithStdErr(errOut),
	)

	// when - stderr로 출력하는 명령어 실행
	err := cmd.RunShell("echo 'error message' >&2")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if out.String() != "" {
		t.Errorf("expected empty stdout, got %s", out.String())
	}

	if !strings.Contains(errOut.String(), "error message") {
		t.Errorf("expected stderr to contain 'error message', got %s", errOut.String())
	}
}

func TestStdIn(t *testing.T) {
	// given
	input := strings.NewReader("hello from stdin\n")
	out := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdIn(input),
		easycmd.WithStdOut(out),
	)

	// when - stdin에서 입력을 읽는 명령어 실행
	err := cmd.Run("cat")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if !strings.Contains(out.String(), "hello from stdin") {
		t.Errorf("expected output to contain 'hello from stdin', got %s", out.String())
	}
}

func TestRunDirConfig(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	tempDir := os.TempDir()

	// RunWithDir 메서드를 사용하여 runDir 설정 테스트
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - RunWithDir를 사용하여 특정 디렉토리에서 실행
	err := cmd.RunWithDir("pwd", tempDir)

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	// 절대 경로로 비교
	expectedDir, _ := filepath.Abs(tempDir)
	actualDir, _ := filepath.Abs(result)

	if expectedDir != actualDir {
		t.Errorf("expected %s, got %s", expectedDir, actualDir)
	}
}

func TestMultipleConfigs(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	input := strings.NewReader("test input\n")

	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithStdErr(errOut),
		easycmd.WithStdIn(input),
	)

	// when - 복합 명령어 실행 (stdin 읽고 stdout과 stderr 모두 사용)
	err := cmd.RunShell("cat && echo 'stderr test' >&2")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if !strings.Contains(out.String(), "test input") {
		t.Errorf("expected stdout to contain 'test input', got %s", out.String())
	}

	if !strings.Contains(errOut.String(), "stderr test") {
		t.Errorf("expected stderr to contain 'stderr test', got %s", errOut.String())
	}
}

func TestCommandParsing(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - 복합 명령어 테스트 (Name과 Args가 올바르게 파싱되는지 간접 확인)
	err := cmd.Run("echo hello world test")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "hello world test" {
		t.Errorf("expected 'hello world test', got '%s'", result)
	}
}

func TestCommandParsingWithQuotes(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - 따옴표가 포함된 명령어 테스트
	err := cmd.Run("echo 'hello world'")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

func TestCommandParsingWithDoubleQuotes(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - 이중 따옴표가 포함된 명령어 테스트
	err := cmd.Run(`echo "hello world"`)

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

func TestSingleCommand(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - 인수가 없는 단일 명령어 테스트
	err := cmd.Run("pwd")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result == "" {
		t.Error("expected non-empty output from pwd command")
	}
}

func TestShellCommandWrapping(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - 쉘 특화 문법 테스트 (bash 래핑이 올바르게 작동하는지 확인)
	err := cmd.RunShell("VAR=test; echo $VAR")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "test" {
		t.Errorf("expected 'test', got '%s'", result)
	}
}

func TestPowershellCommandWrapping(t *testing.T) {
	// PowerShell이 설치되어 있지 않은 경우 스킵
	if _, err := os.Stat("/usr/bin/pwsh"); os.IsNotExist(err) {
		if _, err := os.Stat("/usr/local/bin/pwsh"); os.IsNotExist(err) {
			t.Skip("PowerShell not installed, skipping test")
		}
	}

	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - PowerShell 명령어 테스트
	err := cmd.RunPowershell("Write-Output 'Hello PowerShell'")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if !strings.Contains(result, "Hello PowerShell") {
		t.Errorf("expected output to contain 'Hello PowerShell', got '%s'", result)
	}
}

func TestPowershellWithDir(t *testing.T) {
	// PowerShell이 설치되어 있지 않은 경우 스킵
	if _, err := os.Stat("/usr/bin/pwsh"); os.IsNotExist(err) {
		if _, err := os.Stat("/usr/local/bin/pwsh"); os.IsNotExist(err) {
			t.Skip("PowerShell not installed, skipping test")
		}
	}

	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(easycmd.WithStdOut(out))

	// when - PowerShell 명령어를 특정 디렉토리에서 실행
	err := cmd.RunPowershellWithDir("Get-Location", "..")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := out.String()
	if !strings.Contains(result, "seungyeop-lee") {
		t.Errorf("expected path to contain 'seungyeop-lee', got %s", result)
	}
}

func TestDebugMode(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	debugOut := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithDebug(debugOut),
	)

	// when
	err := cmd.Run("echo hello")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	debugResult := debugOut.String()
	commandResult := out.String()

	// 디버그 출력 확인 (DebugOut에서)
	if !strings.Contains(debugResult, "[DEBUG] 파싱된 명령어: echo hello") {
		t.Errorf("expected debug output to contain parsed command, got %s", debugResult)
	}
	if !strings.Contains(debugResult, "[DEBUG] 실행 명령어: echo") {
		t.Errorf("expected debug output to contain command name, got %s", debugResult)
	}
	if !strings.Contains(debugResult, "[DEBUG] 실행 인수: [hello]") {
		t.Errorf("expected debug output to contain command args, got %s", debugResult)
	}
	if !strings.Contains(debugResult, "[DEBUG] 명령어 실행 시작...") {
		t.Errorf("expected debug output to contain start message, got %s", debugResult)
	}
	if !strings.Contains(debugResult, "[DEBUG] 명령어 실행 완료") {
		t.Errorf("expected debug output to contain completion message, got %s", debugResult)
	}

	// 실제 명령어 출력 확인 (StdOut에서)
	if !strings.Contains(commandResult, "hello") {
		t.Errorf("expected actual command output, got %s", commandResult)
	}

	// StdOut에는 디버그 메시지가 없어야 함
	if strings.Contains(commandResult, "[DEBUG]") {
		t.Errorf("expected no debug messages in command output, got %s", commandResult)
	}
}

func TestDebugModeWithDir(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	debugOut := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithDebug(debugOut),
	)

	// when
	err := cmd.RunWithDir("pwd", "..")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	debugResult := debugOut.String()

	// 디렉토리 디버그 출력 확인
	if !strings.Contains(debugResult, "[DEBUG] 실행 디렉토리: ..") {
		t.Errorf("expected debug output to contain run directory, got %s", debugResult)
	}
}

func TestDebugModeShellCommand(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	debugOut := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithDebug(debugOut),
	)

	// when
	err := cmd.RunShell("echo 'shell test'")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	debugResult := debugOut.String()

	// 쉘 명령어 래핑 확인
	if !strings.Contains(debugResult, "[DEBUG] 파싱된 명령어: bash -c echo 'shell test'") {
		t.Errorf("expected debug output to contain shell wrapped command, got %s", debugResult)
	}
	if !strings.Contains(debugResult, "[DEBUG] 실행 명령어: bash") {
		t.Errorf("expected debug output to contain bash as command name, got %s", debugResult)
	}
}

func TestDebugOutputSeparation(t *testing.T) {
	// given
	stdOut := &bytes.Buffer{}
	debugOut := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(stdOut),
		easycmd.WithDebug(debugOut),
	)

	// when
	err := cmd.Run("echo hello")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	stdOutResult := stdOut.String()
	debugOutResult := debugOut.String()

	// 명령어 출력은 StdOut에만 있어야 함
	if !strings.Contains(stdOutResult, "hello") {
		t.Errorf("expected command output in StdOut, got %s", stdOutResult)
	}
	if strings.Contains(stdOutResult, "[DEBUG]") {
		t.Errorf("expected no debug messages in StdOut, got %s", stdOutResult)
	}

	// 디버그 출력은 DebugOut에만 있어야 함
	if !strings.Contains(debugOutResult, "[DEBUG]") {
		t.Errorf("expected debug messages in DebugOut, got %s", debugOutResult)
	}

	// DebugOut에는 실제 명령어 실행 결과 (줄바꿈 포함)가 없어야 함
	// 디버그 메시지에 명령어 문자열이 포함되는 것은 정상
	lines := strings.Split(debugOutResult, "\n")
	for _, line := range lines {
		if line == "hello" {
			t.Errorf("expected no actual command execution output in DebugOut, found line: %s", line)
		}
	}
}

func TestDebugDefaultOutputToStderr(t *testing.T) {
	// given - DebugOut을 명시적으로 설정하지 않음 (기본값 사용)
	out := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithDebug(),
	)

	// when
	err := cmd.Run("echo 'test'")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := out.String()

	// StdOut에는 명령어 출력만 있고 디버그 메시지는 없어야 함 (기본적으로 stderr로 가므로)
	if !strings.Contains(result, "test") {
		t.Errorf("expected command output in StdOut, got %s", result)
	}
	if strings.Contains(result, "[DEBUG]") {
		t.Errorf("expected no debug messages in StdOut when using default DebugOut, got %s", result)
	}
}

func TestWithTimeout(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithTimeoutSeconds(2),
	)

	// when - 빠르게 실행되는 명령어
	err := cmd.Run("echo 'timeout test'")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "timeout test" {
		t.Errorf("expected 'timeout test', got '%s'", result)
	}
}

func TestWithTimeoutExceeded(t *testing.T) {
	// given
	cmd := easycmd.New(
		easycmd.WithTimeoutSeconds(1),
	)

	// when - 타임아웃보다 오래 걸리는 명령어
	err := cmd.Run("sleep 3")

	// then
	if err == nil {
		t.Error("expected timeout error, got nil")
		return
	}

	// 타임아웃 에러 확인
	errMsg := err.Error()
	if !strings.Contains(errMsg, "명령어 실행 타임아웃") {
		t.Errorf("expected timeout error message, got: %s", errMsg)
	}
}

func TestWithTimeoutMillis(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithTimeoutMillis(2000), // 2초
	)

	// when - 빠르게 실행되는 명령어
	err := cmd.Run("echo 'timeout millis test'")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "timeout millis test" {
		t.Errorf("expected 'timeout millis test', got '%s'", result)
	}
}

func TestWithTimeoutVeryShort(t *testing.T) {
	// given
	cmd := easycmd.New(
		easycmd.WithTimeoutMillis(1), // 1밀리초
	)

	// when - 매우 짧은 타임아웃으로 명령어 실행
	err := cmd.Run("sleep 1")

	// then
	if err == nil {
		t.Error("expected timeout error, got nil")
		return
	}

	// 매우 짧은 타임아웃이므로 시작 실패 또는 실행 타임아웃 둘 다 가능
	errMsg := err.Error()
	if !strings.Contains(errMsg, "명령어 시작 실패") && !strings.Contains(errMsg, "명령어 실행 타임아웃") {
		t.Errorf("expected timeout-related error message, got: %s", errMsg)
	}
}

func TestWithEnv(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	testEnv := []string{"TEST_VAR=hello_world", "ANOTHER_VAR=test_value"}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithEnv(testEnv),
	)

	// when - 환경변수를 출력하는 쉘 명령어
	err := cmd.RunShell("echo $TEST_VAR")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "hello_world" {
		t.Errorf("expected 'hello_world', got '%s'", result)
	}
}

func TestWithEnvMultiple(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	testEnv := []string{
		"VAR1=value1",
		"VAR2=value2",
		"VAR3=value3",
	}
	cmd := easycmd.New(
		easycmd.WithStdOut(out),
		easycmd.WithEnv(testEnv),
	)

	// when - 여러 환경변수를 출력하는 쉘 명령어
	err := cmd.RunShell("echo $VAR1 $VAR2 $VAR3")

	// then
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	result := strings.TrimSpace(out.String())
	if result != "value1 value2 value3" {
		t.Errorf("expected 'value1 value2 value3', got '%s'", result)
	}
}
