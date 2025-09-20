package easycmd

import (
	"fmt"
	"io"
	"time"
)

// Logger 디버그 출력을 담당하는 인터페이스
type Logger interface {
	ParsedCommand(command string)
	ExecutionCommand(name string, args []string)
	ExecutionDirectory(dir string)
	ExecutionStart()
	Timeout(timeout time.Duration)
	Environment(envCount int)
	StartFailed(err error)
	ExecutionFailed(err error, isTimeout bool)
	ExecutionCompleted()
}

// DebugLogger 실제 디버그 출력을 수행하는 구현체
type DebugLogger struct {
	out       io.Writer
	startTime time.Time
}

// NewDebugLogger DebugLogger 인스턴스를 생성합니다
func NewDebugLogger(out io.Writer) *DebugLogger {
	return &DebugLogger{out: out}
}

func (d *DebugLogger) ParsedCommand(command string) {
	fmt.Fprintf(d.out, "[DEBUG] 파싱된 명령어: %s\n", command)
}

func (d *DebugLogger) ExecutionCommand(name string, args []string) {
	fmt.Fprintf(d.out, "[DEBUG] 실행 명령어: %s\n", name)
	fmt.Fprintf(d.out, "[DEBUG] 실행 인수: %v\n", args)
}

func (d *DebugLogger) ExecutionDirectory(dir string) {
	if dir != "" {
		fmt.Fprintf(d.out, "[DEBUG] 실행 디렉토리: %s\n", dir)
	}
}

func (d *DebugLogger) ExecutionStart() {
	fmt.Fprintf(d.out, "[DEBUG] 명령어 실행 시작...\n")
	d.startTime = time.Now()
}

func (d *DebugLogger) Timeout(timeout time.Duration) {
	fmt.Fprintf(d.out, "[DEBUG] 타임아웃 설정: %s\n", timeout)
}

func (d *DebugLogger) Environment(envCount int) {
	fmt.Fprintf(d.out, "[DEBUG] 환경변수 설정: %d개\n", envCount)
}

func (d *DebugLogger) StartFailed(err error) {
	fmt.Fprintf(d.out, "[DEBUG] 명령어 시작 실패: %s\n", err)
}

func (d *DebugLogger) ExecutionFailed(err error, isTimeout bool) {
	if isTimeout {
		fmt.Fprintf(d.out, "[DEBUG] 명령어 실행 타임아웃: %v\n", err)
	} else {
		fmt.Fprintf(d.out, "[DEBUG] 명령어 실행 실패: %v\n", err)
	}
}

func (d *DebugLogger) ExecutionCompleted() {
	actualDuration := time.Since(d.startTime)
	fmt.Fprintf(d.out, "[DEBUG] 명령어 실행 완료 (실행 시간: %s)\n", actualDuration)
}

// NoOpLogger 아무것도 하지 않는 로거 구현체 (Null Object Pattern)
type NoOpLogger struct{}

// NewNoOpLogger NoOpLogger 인스턴스를 생성합니다
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

func (n *NoOpLogger) ParsedCommand(command string)                {}
func (n *NoOpLogger) ExecutionCommand(name string, args []string) {}
func (n *NoOpLogger) ExecutionDirectory(dir string)               {}
func (n *NoOpLogger) ExecutionStart()                             {}
func (n *NoOpLogger) Timeout(timeout time.Duration)               {}
func (n *NoOpLogger) Environment(envCount int)                    {}
func (n *NoOpLogger) StartFailed(err error)                       {}
func (n *NoOpLogger) ExecutionFailed(err error, isTimeout bool)   {}
func (n *NoOpLogger) ExecutionCompleted()                         {}
