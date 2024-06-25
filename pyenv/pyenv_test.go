package pyenv

import (
	"fmt"
	"testing"
)

func TestIntegration(t *testing.T) {
	if !DistExists() {
		MacInstall()
	}
	program := `
		print('hello')
		print('world')
	`
	out, e := Execute(program)
	fmt.Printf("%s : %v", out, e)
}

func TestDependencies(t *testing.T) {
	ret, _ := AddDependencies("../requirements.txt")
	t.Logf("ret: %s", ret)
}
