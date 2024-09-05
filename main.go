package main

import (
	"fmt"
	"os"
	"os/exec"
	osuser "os/user"
	"path/filepath"
	"strings"
)

const (
	// ssh options which takes no value
	optsNoValue = "1246AaCfGgKkMNnqsTtVvXxYy"
	// ssh options which takes value
	optsWithValue = "BbcDEeFIiJLlmOoPpRSWw"
)

// awscli command available
var awsAvailable bool

func init() {
	_, err := exec.LookPath("aws")
	awsAvailable = err == nil
	// _, err = exec.LookPath("ssh")
	// _, err2 := exec.LookPath("ssh-keygen")
	// sshAvailable = err == nil && err2 == nil
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
	user   string
	instId string
	// identity file is specified only if the user explicitly specifies it
	ident string
}

// parse arguments to get values for ec2 instance connect
// Returns:
// 1. arguments after parsing
// 2. values for ec2 instance connect
func argparse() ([]string, *awsargs) {
	args := os.Args[1:]

	// if aws command is not available, checking process is not needed anymore
	if !awsAvailable {
		return args, nil
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
				return args, nil
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
					return args, nil
				}
			}
		} else {
			// user@host or host
			atat := strings.IndexRune(arg, '@')
			if atat == -1 {
				u, err := osuser.Current()
				if err != nil {
					return args, nil
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
				return args, &awsargs{user, host, ident}
			} else {
				return args, nil
			}
		}
		i++
	}
	return args, nil
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

func createSubmitIdent(arg *awsargs, tmpDir string) (string, error) {
	priv := filepath.Join(tmpDir, "id_rsa")
	pub := filepath.Join(tmpDir, "id_rsa.pub")
	exec.Command("ssh-keygen", "-t", "rsa", "-N", "", "-f", priv).Run()

	privAbs, err := filepath.Abs(priv)
	if err != nil {
		return "", err
	}
	pubAbs, err := filepath.Abs(pub)
	if err != nil {
		return "", err
	}

	pubUrl := "file://" + pubAbs

	exec.Command(
		"aws",
		"ec2-instance-connect",
		"send-ssh-public-key",
		"--instance-id",
		arg.instId,
		"--instance-os-user",
		arg.user,
		"--ssh-public-key",
		pubUrl,
	).Run()
	return privAbs, nil
}

func main() {
	args, aws := argparse()
	if aws == nil {
		// fallback to original ssh command
		os.Exit(ssh(args))
	}

	tmpDir, err := os.MkdirTemp("", "eish-tmp-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(ssh(args))
	}
	defer os.RemoveAll(tmpDir)

	customOpts := []string{
		"-o",
		"ProxyCommand=aws ec2-instance-connect open-tunnel --instance-id %h",
	}
	if aws.ident != "" {
		privKey, err := createSubmitIdent(aws, tmpDir)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(ssh(args))
		}
		customOpts = append(customOpts, "-i", privKey)
	}
	os.Exit(ssh(append(customOpts, args...)))
}
