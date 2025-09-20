package easycmd

import (
	"errors"
	"fmt"
	"os/exec"
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
	config := c.c
	config.RunDir = runDir(runDirStr)
	return run(command(commandStr), config)
}

func (c *Cmd) RunShellWithDir(commandStr string, runDirStr string) error {
	config := c.c
	config.RunDir = runDir(runDirStr)
	return run(command(commandStr).ShellCommand(), config)
}

func (c *Cmd) RunPowershellWithDir(commandStr string, runDirStr string) error {
	config := c.c
	config.RunDir = runDir(runDirStr)
	return run(command(commandStr).PowershellCommand(), config)
}

func run(command command, config config) error {
	if command == "" {
		return EmptyCmdError
	}

	if config.Debug {
		fmt.Fprintf(config.DebugOut, "[DEBUG] 파싱된 명령어: %s\n", command.String())
		fmt.Fprintf(config.DebugOut, "[DEBUG] 실행 명령어: %s\n", command.Name())
		fmt.Fprintf(config.DebugOut, "[DEBUG] 실행 인수: %v\n", command.Args())
		if config.RunDir != "" {
			fmt.Fprintf(config.DebugOut, "[DEBUG] 실행 디렉토리: %s\n", string(config.RunDir))
		}
		fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 실행 시작...\n")
	}

	cmd := exec.Command(command.Name(), command.Args()...)

	cmd.Dir = string(config.RunDir)
	cmd.Stdin = config.StdIn
	cmd.Stdout = config.StdOut
	cmd.Stderr = config.StdErr

	if err := cmd.Start(); err != nil {
		if config.Debug {
			fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 시작 실패: %s\n", err)
		}
		return fmt.Errorf("can't start command: %s", err)
	}
	err := cmd.Wait()

	if err != nil {
		if config.Debug {
			fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 실행 실패: %v\n", err)
		}
		return fmt.Errorf("command fails to run or doesn't complete successfully: %v", err)
	}

	if config.Debug {
		fmt.Fprintf(config.DebugOut, "[DEBUG] 명령어 실행 완료\n")
	}

	return nil
}

var EmptyCmdError = errors.New("empty command")
