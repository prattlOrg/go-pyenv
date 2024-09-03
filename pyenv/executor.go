package pyenv

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// distributions: windows_x86, windows_x64, apple_aarch64, apple_x64, linux_gnu_aarch64, linux_gnu_x64, linux_gnu_x64_v2, linux_gnu_x64_v3, linux_gnu_x64_v4
type PyEnv struct {
	ParentPath   string
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
	fp := filepath.Join(env.ParentPath, fmt.Sprintf("dist/python_%s/python/install/bin/pip", env.Distribution))
	cmd := exec.Command(fp, "install", "-r", requirementsPath)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (env *PyEnv) ExecutePython(args ...string) *exec.Cmd {
	pythonCmd := filepath.Join(env.ParentPath, fmt.Sprintf("dist/python_%s/python/install/bin/pip", env.Distribution))
	cmd := exec.Command(pythonCmd, args...)
	return cmd
}
