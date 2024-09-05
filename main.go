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

func hostIsEc2Instance(host string) bool {
	// `^i-([\u\l\d]{8}|[\u\l\d]{17})^`
	if !strings.HasPrefix(host, "i-") {
		return false
	}
	hostr := []rune(host)
	hostlen := len(hostr)
	if hostlen != 10 && hostlen != 19 {
		return false
	}
	for i := 2; i < hostlen; i++ {
		c := hostr[i]
		if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// values for ec2 instance connect
type awsargs struct {
	user  string
	host  string
	ident string
}

// parse arguments to get values for ec2 instance connect
func argparse(args []string) *awsargs {
	// if aws command is not available, checking process is not needed anymore
	if !awsAvailable {
		return nil
	}

	user, host, ident := "", "", ""

	acnt := len(args)
	i := 0
argloop:
	for i < acnt {
		arg := args[i]

		if arg[0] == '-' {
			if len(arg) == 1 {
				// arg is just '-' (invalid)
				return nil
			}
			argr := []rune(arg[1:])
			arglen := len(argr)
		charloop:
			for ci, c := range argr {
				if strings.ContainsRune(optsNoValue, c) {
					continue charloop
				} else if strings.ContainsRune(optsWithValue, c) {
					if ci == arglen-1 {
						i++ // skip the next arg as its value

						// `-l value`
						if c == 'l' {
							user = args[i]
						} else if c == 'i' {
							ident = args[i]
						}
					} else {
						// `-lvalue`
						if c == 'l' {
							user = string(argr[ci+1:])
						} else if c == 'i' {
							ident = string(argr[ci+1:])
						}
					}
					i++ // when continue argloop, i should be increased
					continue argloop
				} else {
					// if the option is not in the list, it should be invalid or unknown option
					return nil
				}
			}
		} else {
			// user@host or host
			atat := strings.IndexRune(arg, '@')
			if atat == -1 {
				u, err := osuser.Current()
				if err != nil {
					return nil
				}
				user = u.Username
				host = arg
			} else {
				host = arg[atat+1:]
				if user == "" {
					user = arg[:atat]
				}
			}

			// host is ec2 instance id pattern or not
			if hostIsEc2Instance(host) {
				return &awsargs{user, host, ident}
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

	aws := argparse(args)
	if aws == nil {
		// fallback to original ssh command
		os.Exit(ssh(args))
	}

	// Debug
	fmt.Printf("%s@%s\n", aws.user, aws.host)

	os.Exit(ssh(args))
}
