package easycmd

import (
	"context"
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
		RunDir:  runDir(runDirStr),
		StdIn:   original.StdIn,
		StdOut:  original.StdOut,
		StdErr:  original.StdErr,
		Logger:  original.Logger,
		Timeout: original.Timeout,
		Env:     original.Env,
	}
}

func run(command command, config config) error {
	if command == "" {
		return EmptyCmdError
	}

	config.Logger.ParsedCommand(command.String())
	config.Logger.ExecutionCommand(command.Name(), command.Args())
	config.Logger.ExecutionDirectory(string(config.RunDir))
	config.Logger.ExecutionStart()

	var cmd *exec.Cmd
	var ctx context.Context
	var cancel context.CancelFunc

	if config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, command.Name(), command.Args()...)
		config.Logger.Timeout(config.Timeout)
	} else {
		cmd = exec.Command(command.Name(), command.Args()...)
	}

	cmd.Dir = string(config.RunDir)
	cmd.Stdin = config.StdIn
	cmd.Stdout = config.StdOut
	cmd.Stderr = config.StdErr
	if len(config.Env) > 0 {
		cmd.Env = config.Env
		config.Logger.Environment(len(config.Env))
	}

	if err := cmd.Start(); err != nil {
		config.Logger.StartFailed(err)
		return fmt.Errorf("명령어를 시작할 수 없습니다: %s", err)
	}
	err := cmd.Wait()

	if err != nil {
		isTimeout := ctx != nil && errors.Is(ctx.Err(), context.DeadlineExceeded)
		config.Logger.ExecutionFailed(err, isTimeout)
		if isTimeout {
			return fmt.Errorf("명령어 실행이 타임아웃되었습니다 (%s): %v", config.Timeout, err)
		}
		return fmt.Errorf("명령어 실행이 실패했거나 성공적으로 완료되지 않았습니다: %v", err)
	}

	config.Logger.ExecutionCompleted()

	return nil
}

var EmptyCmdError = errors.New("empty command")
