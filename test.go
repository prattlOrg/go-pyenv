package test

import (
	"fmt"
	"testing"

	"github.com/voidKandy/go-pyenv/pyenv"
)

func TestIntegration(t *testing.T) {
	if !pyenv.DistExists() {
		pyenv.MacInstall()
	}
	program := `
		print('hello')
		print('world')
	`
	out, e := pyenv.Execute(program)
	fmt.Printf("%s : %v", out, *e)
}
