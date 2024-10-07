package pyenv

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func compareDirectories(dir1, dir2 string) (bool, error) {
	files1, err := os.ReadDir(dir1)
	if err != nil {
		return false, err
	}
	files2, err := os.ReadDir(dir2)
	if err != nil {
		return false, err
	}

	if len(files1) != len(files2) {
		return false, nil
	}

	fileMap := make(map[string]os.FileInfo)
	for _, file := range files2 {
		info, err := file.Info()
		if err != nil {
			return false, err
		}
		fileMap[file.Name()] = info
	}

	for _, file1 := range files1 {
		file2info, exists := fileMap[file1.Name()]
		file1info, err := file1.Info()
		if err != nil {
			return false, err
		}
		if !exists || file1info.Size() != file2info.Size() {
			return false, nil
		}
	}

	return true, nil
}

func TestCompression(t *testing.T) {
	dir := "../comp-test"
	zipTarget := "../test.zip"
	os.ReadDir(dir)

	err := compressDir(dir, zipTarget)
	if err != nil {
		log.Fatalf("Error unzipping: %v\n", err)
	}
	unzipTarget := "../target"
	err = unzipSource(zipTarget, unzipTarget)
	if err != nil {
		log.Fatalf("Error unzipping: %v\n", err)
	}

	same, err := compareDirectories(dir, unzipTarget)
	if err != nil {
		log.Fatalf("Error unzipping: %v\n", err)
	}

	if !same {
		log.Fatalf("Compression/Decompression didn't, fail, but it didn't output the expected directory contents")
	}

	os.RemoveAll(unzipTarget)
	os.RemoveAll(zipTarget)
}

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
	cmd, err := env.ExecutePython("c", program)
	if err != nil {
		t.Errorf("error executing python: %v\n", err)
	}

	cmdT := fmt.Sprintf("%T", cmd)
	t.Log(cmdT)

	err = env.CompressDist()
	if err != nil {
		t.Errorf("error compressing dist: %v\n", err)
	}

	if _, err = env.ExecutePython("c", program); err == nil {
		t.Error("Execute python should error when trying to run when dist is compressed")
	}
	t.Log("compressed & Execute python returned as expected")

	if err := env.DecompressDist(); err != nil {
		t.Errorf("error decompressing dist: %v\n", err)
	}

	t.Log("decompressed")

	cmd, err = env.ExecutePython("c", program)
	if err != nil {
		t.Errorf("error executing python: %v\n", err)
	}

	cmdT2 := fmt.Sprintf("%T", cmd)
	if cmdT != cmdT2 {
		t.Logf("expected outputs to be the same. Instead got cmdT: %v\ncmdT2: %v\n", cmdT, cmdT2)
	}

	t.Log("Test passed")
}

func TestDependencies(t *testing.T) {
	env := testEnv()
	_ = env.AddDependencies("requirements.txt")
	list, _ := env.executePip("list")
	// t.Logf("ret: %s", ret)
	t.Logf("list: %s", list)
}

func TestRemove(t *testing.T) {
	env := testEnv()
	err := os.RemoveAll(env.ParentPath)
	if err != nil {
		t.Logf("Problem cleaning %s: %v", env.ParentPath, err)
	}
	t.Log("Successfully cleaned prattl directory")
}

func (env *PyEnv) executePip(arg string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdPath := filepath.Join(env.ParentPath, "dist/python/install/bin/pip")
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
	env, _ := NewPyEnv(filepath.Join(dirname, ".pyenv_test"))
	env.Distribution = "darwin/arm64"
	return *env
}
