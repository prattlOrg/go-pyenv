package pyenv

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

func DistExists() bool {
	return fs.ValidPath("dist")
}

func AddDependencies(requirementsPath string) (string, error) {
	deps, err := os.ReadFile(requirementsPath)
	if err != nil {
		return "", err
	}

	fmt.Println(string(deps))

	// var out bytes.Buffer
	// var stderr bytes.Buffer
	// cmd := exec.Command("./dist/python-mac.extracted/python/install/bin/pip", "install", )
	// cmd.Stdout = &out
	// cmd.Stderr = &stderr
	// if err := cmd.Start(); err != nil {
	// 	e := fmt.Errorf(stderr.String())
	// 	return "", &e
	// }
	// if err := cmd.Wait(); err != nil {
	// 	e := fmt.Errorf(stderr.String())
	// 	return "", &e
	// }
	// e := fmt.Errorf(stderr.String())
	// output := out.String()
	// return output, &e
	return "", nil
}

func Execute(arg string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("./dist/python-mac.extracted/python/install/bin/python", "-c", arg)
	// cmd := exec.Command("python", "-v")
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		e := fmt.Errorf(stderr.String())
		return "", e
	}
	if err := cmd.Wait(); err != nil {
		e := fmt.Errorf(stderr.String())
		return "", e
	}
	e := fmt.Errorf(stderr.String())
	output := out.String()
	return output, e
}
