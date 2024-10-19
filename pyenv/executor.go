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
		return nil, fmt.Errorf("error getting $HOME directory: %v", err)
	}

	if path == homedir {
		err := fmt.Errorf("path cannot be homedir\npath given: %s\nhomedir: %s", path, homedir)
		return nil, err
	}

	env := PyEnv{
		ParentPath: path,
	}

	return &env, nil
}

func (env *PyEnv) distPath() string {
	return filepath.Join(env.ParentPath, "dist")
}

func (env *PyEnv) compressionTarget() string {
	return env.distPath() + ZIP_FILE_EXT
}

func (env *PyEnv) CompressDist() error {
	if env.Compressed {
		return fmt.Errorf("dist is already compressed")
	}

	if err := compressDir(env.distPath(), env.compressionTarget()); err != nil {
		return fmt.Errorf("error compressing python environment: %v", err)
	}
	env.Compressed = true

	if err := os.RemoveAll(env.distPath()); err != nil {
		return fmt.Errorf("error removing old uncompressed evironment: %v", err)
	}
	log.Printf("removed %v\n", env.distPath())
	return nil
}

func (env *PyEnv) DecompressDist() error {
	if !env.Compressed {
		log.Println("dist is already decompressed")
		return nil
	}

	env.Compressed = false

	if err := unzipSource(env.compressionTarget(), env.distPath()); err != nil {
		return fmt.Errorf("error unzipping compressed evironment: %v", err)
	}
	if err := os.RemoveAll(env.compressionTarget()); err != nil {
		return fmt.Errorf("error removing old compressed evironment: %v", err)
	}
	log.Printf("removed %v\n", env.compressionTarget())
	return nil
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
	if env.Compressed {
		if err := env.DecompressDist(); err != nil {
			return err
		}
	}
	var fp string
	if strings.Contains(env.Distribution, "windows") {
		fp = filepath.Join(env.ParentPath, "dist/python/install/Scripts/pip3.exe")
	} else {
		fp = filepath.Join(env.ParentPath, "dist/python/install/bin/pip")
	}
	log.Println("installing python dependencies")
	cmd := exec.Command(fp, "install", "-r", requirementsPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error installing python dependencies: %v", err)
	}
	log.Println("installing python dependencies complete")
	return nil
}

// Executes given python arguments, only will error if the env is compressed
func (env *PyEnv) ExecutePython(args ...string) (*exec.Cmd, error) {
	if env.Compressed {
		return nil, fmt.Errorf("cannot execute python with a compressed dist")
	}
	var fp string
	if strings.Contains(env.Distribution, "windows") {
		fp = filepath.Join(env.ParentPath, "dist/python/install/python.exe")
	} else {
		fp = filepath.Join(env.ParentPath, "dist/python/install/bin/python")
	}
	cmd := exec.Command(fp, args...)
	return cmd, nil
}
