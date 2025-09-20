package easycmd

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

type Cmd struct {
	c config
}

type configApply func(c *config)

func New(configApplies ...configApply) *Cmd {
	c := &config{}
	for _, ca := range configApplies {
		ca(c)
	}
	c.fillDefault()

	return &Cmd{
		c: *c,
	}
}

func (c *Cmd) Run(commandStr string) error {
	return run(command(commandStr), c.c)
}

func (c *Cmd) RunShell(commandStr string) error {
	return run(command(commandStr).ShellCommand(), c.c)
}

func (c *Cmd) RunPowershell(commandStr string) error {
	return run(command(commandStr).PowershellCommand(), c.c)
}

func (c *Cmd) RunWithDir(commandStr string, runDirStr string) error {
	config := copyConfigWithDir(c.c, runDirStr)
	return run(command(commandStr), config)
}

func (c *Cmd) RunShellWithDir(commandStr string, runDirStr string) error {
	config := copyConfigWithDir(c.c, runDirStr)
	return run(command(commandStr).ShellCommand(), config)
}

func (c *Cmd) RunPowershellWithDir(commandStr string, runDirStr string) error {
	config := copyConfigWithDir(c.c, runDirStr)
	return run(command(commandStr).PowershellCommand(), config)
}

func copyConfigWithDir(original config, runDirStr string) config {
	return config{
		RunDir:   runDir(runDirStr),
		StdIn:    original.StdIn,
		StdOut:   original.StdOut,
		StdErr:   original.StdErr,
		Debug:    original.Debug,
		DebugOut: original.DebugOut,
		Timeout:  original.Timeout,
		Env:      original.Env,
	}
}

func run(command command, config config) error {
	if command == "" {
		return EmptyCmdError
	}

	var startTime time.Time
	if config.Debug {
		fmt.Fprintf(config.DebugOut, "[DEBUG] 파싱된 명령어: %s\n", command.String())
		fmt.Fprintf(config.DebugOut, "[DEBUG] 실행 명령어: %s\n", command.Name())
		fmt.Fprintf(config.DebugOut, "[DEBUG] 실행 인수: %v\n", command.Args())
		if config.RunDir != "" {
			fmt.Fprintf(config.DebugOut, "[DEBUG] 실행 디렉토리: %s\n", string(config.RunDir))
		}
		fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 실행 시작...\n")
		startTime = time.Now()
	}

	var cmd *exec.Cmd
	var ctx context.Context
	var cancel context.CancelFunc

	if config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, command.Name(), command.Args()...)
		if config.Debug {
			fmt.Fprintf(config.DebugOut, "[DEBUG] 타임아웃 설정: %s\n", config.Timeout)
		}
	} else {
		cmd = exec.Command(command.Name(), command.Args()...)
	}

	cmd.Dir = string(config.RunDir)
	cmd.Stdin = config.StdIn
	cmd.Stdout = config.StdOut
	cmd.Stderr = config.StdErr
	if len(config.Env) > 0 {
		cmd.Env = config.Env
		if config.Debug {
			fmt.Fprintf(config.DebugOut, "[DEBUG] 환경변수 설정: %d개\n", len(config.Env))
		}
	}

	if err := cmd.Start(); err != nil {
		if config.Debug {
			fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 시작 실패: %s\n", err)
		}
		return fmt.Errorf("명령어를 시작할 수 없습니다: %s", err)
	}
	err := cmd.Wait()

	if err != nil {
		if config.Debug {
			if ctx != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
				fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 실행 타임아웃: %v\n", err)
			} else {
				fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 실행 실패: %v\n", err)
			}
		}
		if ctx != nil && errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("명령어 실행이 타임아웃되었습니다 (%s): %v", config.Timeout, err)
		}
		return fmt.Errorf("명령어 실행이 실패했거나 성공적으로 완료되지 않았습니다: %v", err)
	}

	if config.Debug {
		duration := time.Since(startTime)
		fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 실행 완료 (실행 시간: %s)\n", duration)
	}

	return nil
}

var EmptyCmdError = errors.New("empty command")
