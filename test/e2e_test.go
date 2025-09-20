package test

import (
	"bytes"
	"testing"

	"github.com/seungyeop-lee/easycmd"
)

func TestSimple(t *testing.T) {
	// given
	out := &bytes.Buffer{}
	cmd := easycmd.New(func(c *easycmd.Config) {
		c.StdOut = out
	})

	// when
	err := cmd.Run("echo hello world")

	// then
	if out.String() != "hello world\n" {
		t.Errorf("expected hello world, got %s", out.String())
	}
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRunShell(t *testing.T) {
	cmd := easycmd.New()

	err := cmd.RunShell("(cd .. && pwd)")

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRunMultiLineShell(t *testing.T) {
	cmd := easycmd.New()

	err := cmd.RunShell(`
	pwd
	ls -al
`)

	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
