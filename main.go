package main

import (
	"fmt"
	"io/fs"

	"github.com/voidKandy/go-pyenv/pyenv"
)

func main() {
	if !fs.ValidPath("dist") {
		pyenv.MacInstall()
	}

	program := `
		print('hello')
		print('world')
	`
	out, e := pyenv.Execute(program)
	fmt.Printf("%s : %v", out, *e)
}
