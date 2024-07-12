package pyenv

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

func (env *PyEnv) AddDependencies(requirementsPath string) error {
	fp := filepath.Join(env.ParentPath, "dist/python-mac.extracted/python/install/bin/pip")
	cmd := exec.Command(fp, "install", "-r", requirementsPath)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (env *PyEnv) ExecutePython(args ...string) *exec.Cmd {
	cmd := exec.Command(env.ParentPath+"dist/python-mac.extracted/python/install/bin/python", args...)
	return cmd
}
