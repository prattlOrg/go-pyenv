package pyenv

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIntegration(t *testing.T) {
	env := testEnv()
	exists, _ := env.DistExists()
	if !*exists {
		env.Install()
	}
	program := `
print('hello')
print('world')
	`
	cmd := env.ExecutePython("c", program)
	cmdT := fmt.Sprintf("%T", cmd)
	fmt.Println(cmdT)
}

func TestDependencies(t *testing.T) {
	env := testEnv()
	_ = env.AddDependencies("requirements.txt")
	list, _ := env.executePip("list")
	// t.Logf("ret: %s", ret)
	t.Logf("list: %s", list)
}

func (env *PyEnv) executePip(arg string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdPath := filepath.Join(env.ParentPath, fmt.Sprintf("dist/python_%s/python/install/bin/pip", env.Distribution))
	cmd := exec.Command(cmdPath, arg)
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

func testEnv() PyEnv {
	dirname, _ := os.UserHomeDir()
	return PyEnv{
		ParentPath:   filepath.Join(dirname, ".pyenv_test"),
		Distribution: "windows_x64",
	}
}
