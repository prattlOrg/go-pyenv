package pyenv

import (
	"fmt"
	// "os/exec"
	"testing"
)

func testEnv() PyEnv {
	return PyEnv{
		ParentPath: "../",
	}
}

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
	out, e := env.ExecutePython(program)
	fmt.Printf("%s : %v", out, e)
}

func TestDependencies(t *testing.T) {
	env := testEnv()
	ret, _ := env.AddDependencies("../requirements.txt")
	list, _ := env.executePip([]string{"list"})
	t.Logf("ret: %s", ret)
	t.Logf("list: %s", list)
}
