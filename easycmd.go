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
	return strings.Split(string(c), " ")[0]
}

func (c command) Args() []string {
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

func (c command) String() string {
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

func (c *Cmd) Run(commandStr string) error {
	return run(command(commandStr), c.c.RunDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunShell(commandStr string) error {
	return run(command(commandStr).ShellCommand(), c.c.RunDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunPowershell(commandStr string) error {
	return run(command(commandStr).PowershellCommand(), c.c.RunDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunWithDir(commandStr string, runDir string) error {
	return run(command(commandStr), RunDir(runDir), c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunShellWithDir(commandStr string, runDir string) error {
	return run(command(commandStr).ShellCommand(), RunDir(runDir), c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func (c *Cmd) RunPowershellWithDir(commandStr string, runDir string) error {
	return run(command(commandStr).PowershellCommand(), RunDir(runDir), c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func run(command command, runDir RunDir, stdin StdIn, stdout StdOut, stderr StdErr) error {
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
