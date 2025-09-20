package easycmd

import "strings"

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
