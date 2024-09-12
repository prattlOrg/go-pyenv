package pyenv

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PyEnv struct {
	ParentPath string
	// distributions: windows/386 windows/amd64 darwin/amd64 darwin/arm64 linux/arm64 linux/amd64
	Distribution string
}

func DefaultPyEnv() PyEnv {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return PyEnv{
		ParentPath: dirname,
	}
}

func (env *PyEnv) DistExists() (*bool, error) {
	fp := filepath.Join(env.ParentPath, "dist")
	_, err := os.Stat(fp)
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
	var fp string
	if strings.Contains(env.Distribution, "windows") {
		fp = filepath.Join(env.ParentPath, "dist/python/install/Scripts/pip3.exe")
	} else {
		fp = filepath.Join(env.ParentPath, "dist/python/install/bin/pip")
	}
	cmd := exec.Command(fp, "install", "-r", requirementsPath)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (env *PyEnv) ExecutePython(args ...string) *exec.Cmd {
	var fp string
	if strings.Contains(env.Distribution, "windows") {
		fp = filepath.Join(env.ParentPath, "dist/python/install/python.exe")
	} else {
		fp = filepath.Join(env.ParentPath, "dist/python/install/bin/python")
	}
	cmd := exec.Command(fp, args...)
	return cmd
}
