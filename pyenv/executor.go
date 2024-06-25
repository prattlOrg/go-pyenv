package pyenv

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type PyEnv struct {
	ParentPath string
}

func DefaultPyEnv() PyEnv {
	return PyEnv{
		ParentPath: "./",
	}
}

func (env *PyEnv) DistExists() (*bool, error) {
	_, err := os.Stat(env.ParentPath + "dist")
	t := true
	f := false
	if err == nil {
		return &t, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return &f, nil
	}
	return nil, err
}

func (env *PyEnv) AddDependencies(requirementsPath string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(env.ParentPath+"dist/python-mac.extracted/python/install/bin/pip",
		"install", "-r", requirementsPath,
	)
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

func (env *PyEnv) ExecutePython(arg string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(env.ParentPath+"dist/python-mac.extracted/python/install/bin/python", "-c", arg)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		log.Printf("Error with command start: %s\n", stderr.String())
		e := fmt.Errorf(stderr.String())
		return "", e
	}
	if err := cmd.Wait(); err != nil {
		log.Printf("Error with command wait: %s\n", stderr.String())
		e := fmt.Errorf(stderr.String())
		return "", e
	}
	// e := fmt.Errorf(stderr.String())
	// log.Printf("Error: %s\n", stderr.String())
	output := out.String()
	return output, nil
}
