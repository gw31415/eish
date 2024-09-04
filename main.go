package main

import (
	"fmt"
	"os"
	"os/exec"
	osuser "os/user"
	"strings"
)

const (
	// ssh options which takes no value
	optsNoValue = "1246AaCfGgKkMNnqsTtVvXxYy"
	// ssh options which takes value
	optsWithValue = "BbcDEeFIiJLlmOoPpRSWw"
)

var (
	// awscli command available
	awsAvailable bool
	// ssh command available
	sshAvailable bool
)

func init() {
	_, err := exec.LookPath("aws")
	awsAvailable = err == nil

	_, err = exec.LookPath("ssh")
	sshAvailable = err == nil
}

// values for ec2 instance connect
type addr struct {
	user string
	host string
}

// parse arguments to get values for ec2 instance connect
func argparse(args []string) *addr {
	// if aws command is not available, checking process is not needed anymore
	if !awsAvailable {
		return nil
	}

	user, host := "", ""

	acnt := len(args)
	i := 0
	for i < acnt {
		arg := args[i]

		if arg[0] == '-' {
			if len(arg) == 1 {
				// arg is just '-' (invalid)
				return nil
			} else if len(arg) == 2 {
				// type: -x value
				opt := rune(arg[1])
				if strings.ContainsRune(optsWithValue, opt) {
					// skip the next arg as its value
					i++
					if opt == 'l' {
						user = args[i]
					}
				} else if !strings.ContainsRune(optsNoValue, opt) {
					// if the option is not in the list, it should be invalid
					return nil
				}
			} else {
				// type: -xvalue
				opt := rune(arg[1])
				if opt == 'l' {
					user = arg[2:]
				} else if !strings.ContainsRune(optsWithValue, opt) {
					// if the option is not in the list, it should be invalid
					return nil
				}
			}
		} else {
			// user@host or host
			if user != "" {
				host = arg
			} else {
				// if user is not specified, it should be user@host or user is same as host
				atat := strings.IndexRune(arg, '@')
				if atat == -1 {
					u, err := osuser.Current()
					if err != nil {
						return nil
					}
					user = u.Username
					host = arg
				} else {
					user = arg[:atat]
					host = arg[atat+1:]
				}
			}

			// host is ec2 instance id pattern or not
			if strings.HasPrefix(host, "i-") {
				return &addr{user, host}
			} else {
				return nil
			}
		}
		i++
	}
	return nil
}

// launch ssh command
func ssh(args []string) int {
	ssh := exec.Command("ssh", args...)

	ssh.Stdout = os.Stdout
	ssh.Stderr = os.Stderr
	ssh.Stdin = os.Stdin

	ssh.Run()
	return ssh.ProcessState.ExitCode()
}

func main() {
	args := os.Args[1:]

	adr := argparse(args)
	if adr == nil {
		os.Exit(ssh(args))
	}

	// Debug
	fmt.Printf("%s@%s\n", adr.user, adr.host)

	os.Exit(ssh(args))
}
