package pyenv

import (
	"errors"
	"fmt"
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
	Compressed   bool
}

func NewPyEnv(path string) (*PyEnv, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if path == homedir {
		err := fmt.Errorf("path cannot be homedir\npath given: %s\nhomedir: %s\n", path, homedir)
		return nil, err
	}

	env := PyEnv{
		ParentPath: path,
	}

	return &env, nil
}

// func DefaultPyEnv() PyEnv {
// 	dirname, err := os.UserHomeDir()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return PyEnv{
// 		ParentPath: dirname,
// 	}
// }

// func (env *PyEnv) CompressDist() error {
// 	if env.Compressed {
// 		log.Println("dist is already compressed")
// 		return nil
// 	}
//
// 	source := filepath.Join(env.ParentPath, "dist")
// 	compressDir(source)
//
// 	return nil
// }
//
// func (env *PyEnv) DecompressDist() error {
// 	if !env.Compressed {
// 		log.Println("dist is already decompressed")
// 		return nil
// 	}
// 	fp := filepath.Join(env.ParentPath, "dist")
// 	return nil
// }

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
	log.Println("installing python dependencies")
	cmd := exec.Command(fp, "install", "-r", requirementsPath)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	log.Println("installing python dependencies complete")
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
