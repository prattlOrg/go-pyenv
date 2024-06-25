package pyenv

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Execute(arg string) (string, *error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("./dist/python-mac.extracted/python/install/bin/python", "-c", arg)
	// cmd := exec.Command("python", "-v")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		e := fmt.Errorf(stderr.String())
		return "", &e
	}
	if err := cmd.Wait(); err != nil {
		e := fmt.Errorf(stderr.String())
		return "", &e
	}
	e := fmt.Errorf(stderr.String())
	output := out.String()
	return output, &e
}
