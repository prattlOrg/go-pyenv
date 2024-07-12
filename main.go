package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/voidKandy/go-pyenv/pyenv"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	parentPath := filepath.Join(home, ".prattl/")
	env := pyenv.PyEnv{
		ParentPath: parentPath,
	}
	exists, _ := env.DistExists()
	if !*exists {
		env.MacInstall()
		dependenciesPath := filepath.Join(parentPath, "requirements.txt")
		err = env.AddDependencies(dependenciesPath)
		if err != nil {
			fmt.Println("AddDependencies went wrong")
		}
	}
}
