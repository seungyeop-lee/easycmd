package easycmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type command string

var bashPrefix command = "bash -c "
var powershellPrefix command = "powershell.exe "

func (c command) ShellCommand() command {
	return bashPrefix + c
}

func (c command) PowershellCommand() command {
	return powershellPrefix + c
}

func (c command) Name() string {
	args := parseCommandArgs(string(c))
	if len(args) == 0 {
		return ""
	}
	return args[0]
}

func (c command) Args() []string {
	command := string(c)
	switch true {
	case strings.HasPrefix(command, bashPrefix.String()):
		return []string{"-c", strings.ReplaceAll(command, bashPrefix.String(), "")}
	case strings.HasPrefix(command, powershellPrefix.String()):
		return []string{strings.ReplaceAll(command, powershellPrefix.String(), "")}
	default:
		args := parseCommandArgs(command)
		if len(args) <= 1 {
			return []string{}
		}
		return args[1:]
	}
}

func (c command) String() string {
	return string(c)
}

type runDir string
type stdIn io.Reader
type stdOut io.Writer
type stdErr io.Writer

type config struct {
	RunDir   runDir
	StdIn    stdIn
	StdOut   stdOut
	StdErr   stdErr
	Debug    bool
	DebugOut stdOut
}

func (c *config) fillDefault() {
	if c.StdIn == nil {
		c.StdIn = os.Stdin
	}
	if c.StdOut == nil {
		c.StdOut = os.Stdout
	}
	if c.StdErr == nil {
		c.StdErr = os.Stderr
	}
	if c.DebugOut == nil {
		c.DebugOut = os.Stderr
	}
}

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

func WithDebug(debugOut ...io.Writer) configApply {
	return func(c *config) {
		c.Debug = true
		if len(debugOut) > 0 {
			c.DebugOut = debugOut[0]
		}
	}
}

func WithStdIn(reader io.Reader) configApply {
	return func(c *config) {
		c.StdIn = reader
	}
}

func WithStdOut(writer io.Writer) configApply {
	return func(c *config) {
		c.StdOut = writer
	}
}

func WithStdErr(writer io.Writer) configApply {
	return func(c *config) {
		c.StdErr = writer
	}
}
