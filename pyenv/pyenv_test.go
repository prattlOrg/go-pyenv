package pyenv

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func TestIntegration(t *testing.T) {
	env := testEnv()
	exists, _ := env.DistExists()
	if !*exists {
		env.MacInstall()
	}
	program := `
print('hello')
print('world')
	`
	cmd := env.ExecutePython("c", program)
	cmdT := fmt.Sprintf("%T", cmd)
	t.Log(cmdT)
}

func TestDependencies(t *testing.T) {
	env := testEnv()
	ret, _ := env.AddDependencies("./requirements.txt")
	list, _ := env.executePip("list")
	t.Logf("ret: %s", ret)
	t.Logf("list: %s", list)
}

func (env *PyEnv) executePip(arg string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(env.ParentPath+"dist/python-mac.extracted/python/install/bin/pip",
		arg)

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
	return PyEnv{
		ParentPath: "../",
	}
}
