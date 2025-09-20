package easycmd

import (
	"io"
	"os"
)

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
