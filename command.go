package easycmd

import "strings"

type command string

const (
	bashPrefixStr       = "bash -c "
	powershellPrefixStr = "powershell.exe "
)

var bashPrefix command = bashPrefixStr
var powershellPrefix command = powershellPrefixStr

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
	if strings.HasPrefix(command, bashPrefix.String()) {
		return []string{"-c", strings.TrimPrefix(command, bashPrefix.String())}
	}
	if strings.HasPrefix(command, powershellPrefix.String()) {
		return []string{"-Command", strings.TrimPrefix(command, powershellPrefix.String())}
	}

	args := parseCommandArgs(command)
	if len(args) <= 1 {
		return []string{}
	}
	return args[1:]
}

func (c command) String() string {
	return string(c)
}
