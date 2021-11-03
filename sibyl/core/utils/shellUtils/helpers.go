package shellUtils

import (
	"bytes"
	"os/exec"
)

func Shellout(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// GetGitStats function will return the git stats in the following format:
// "d8e6e45 \n d8e6e45d52f7bf164a995e22abb81ffc6e3eeae1 \n 3 0"
func GetGitStats() string {
	stdout, stderr, err := Shellout(gitCmd)
	if err != nil || len(stderr) != 0 {
		return ""
	}

	return stdout
}
