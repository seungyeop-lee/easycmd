package easycmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Command string

var bashPrefix Command = "bash -c "
var powershellPrefix Command = "powershell.exe "

func (c Command) ShellCommand() Command {
	return bashPrefix + c
}

func (c Command) PowershellCommand() Command {
	return powershellPrefix + c
}

func (c Command) Name() string {
	return strings.Split(string(c), " ")[0]
}

func (c Command) Args() []string {
	command := string(c)
	switch true {
	case strings.HasPrefix(command, bashPrefix.String()):
		return []string{"-c", strings.ReplaceAll(command, bashPrefix.String(), "")}
	case strings.HasPrefix(command, powershellPrefix.String()):
		return []string{strings.ReplaceAll(command, powershellPrefix.String(), "")}
	default:
		return strings.Split(command, " ")[1:]
	}
}

func (c Command) String() string {
	return string(c)
}

type RunDir string
type StdIn io.Reader
type StdOut io.Writer
type StdErr io.Writer

type Config struct {
	RunDir RunDir
	StdIn  StdIn
	StdOut StdOut
	StdErr StdErr
}

func (c *Config) fillDefault() {
	if c.StdIn == nil {
		c.StdIn = os.Stdin
	}
	if c.StdOut == nil {
		c.StdOut = os.Stdout
	}
	if c.StdErr == nil {
		c.StdErr = os.Stderr
	}
}

type Cmd struct {
	c Config
}

type configApply func(c *Config)

func New(configApplies ...configApply) *Cmd {
	c := &Config{}
	for _, ca := range configApplies {
		ca(c)
	}
	c.fillDefault()

	return &Cmd{
		c: *c,
	}
}

func (c *Cmd) Run(command Command) error {
	return run(command, c.c.RunDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunShell(command Command) error {
	return run(command.ShellCommand(), c.c.RunDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunPowershell(command Command) error {
	return run(command.PowershellCommand(), c.c.RunDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunWithDir(command Command, runDir RunDir) error {
	return run(command, runDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunShellWithDir(command Command, runDir RunDir) error {
	return run(command.ShellCommand(), runDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunPowershellWithDir(command Command, runDir RunDir) error {
	return run(command.PowershellCommand(), runDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func run(command Command, runDir RunDir, stdin StdIn, stdout StdOut, stderr StdErr) error {
	if command == "" {
		return EmptyCmdError
	}

	cmd := exec.Command(command.Name(), command.Args()...)

	cmd.Dir = string(runDir)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("can't start command: %s", err)
	}
	err := cmd.Wait()

	if err != nil {
		return fmt.Errorf("command fails to run or doesn't complete successfully: %v", err)
	}

	return nil
}

var EmptyCmdError = errors.New("empty command")
