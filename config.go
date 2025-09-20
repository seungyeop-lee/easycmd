package easycmd

import (
	"io"
	"os"
	"time"
)

type runDir string
type stdIn io.Reader
type stdOut io.Writer
type stdErr io.Writer

type config struct {
	RunDir  runDir
	StdIn   stdIn
	StdOut  stdOut
	StdErr  stdErr
	Logger  Logger
	Timeout time.Duration
	Env     []string
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
	if c.Logger == nil {
		c.Logger = NewNoOpLogger()
	}
}

func WithDebug(debugOut ...io.Writer) configApply {
	return func(c *config) {
		var out io.Writer = os.Stderr
		if len(debugOut) > 0 {
			out = debugOut[0]
		}
		c.Logger = NewDebugLogger(out)
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

func WithTimeout(timeout time.Duration) configApply {
	return func(c *config) {
		c.Timeout = timeout
	}
}

func WithTimeoutSeconds(seconds int) configApply {
	return func(c *config) {
		c.Timeout = time.Duration(seconds) * time.Second
	}
}

func WithTimeoutMillis(millis int) configApply {
	return func(c *config) {
		c.Timeout = time.Duration(millis) * time.Millisecond
	}
}

func WithEnv(env []string) configApply {
	return func(c *config) {
		c.Env = env
	}
}
