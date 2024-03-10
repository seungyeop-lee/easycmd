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

func (c *Cmd) RunWithDir(command Command, runDir RunDir) error {
	return run(command, runDir, c.c.StdIn, c.c.StdOut, c.c.StdErr)
}

func run(command Command, runDir RunDir, stdin StdIn, stdout StdOut, stderr StdErr) error {
	if command == "" {
		return EmptyCmdError
	}

	args := strings.Split(string(command), " ")
	cmd := exec.Command(args[0], args[1:]...)

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
