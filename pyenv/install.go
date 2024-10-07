package pyenv

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

const DIST_DIR = "dist"

var versions = map[string]string{
	"windows/386":   "cpython-3.12.5+20240814-i686-pc-windows-msvc-pgo-full.tar.zst",
	"windows/amd64": "cpython-3.12.5+20240814-x86_64-pc-windows-msvc-pgo-full.tar.zst",
	"darwin/arm64":  "cpython-3.12.5+20240814-aarch64-apple-darwin-pgo+lto-full.tar.zst",
	"darwin/amd64":  "cpython-3.12.5+20240814-x86_64-apple-darwin-pgo+lto-full.tar.zst",
	"linux/arm64":   "cpython-3.12.5+20240814-aarch64-unknown-linux-gnu-lto-full.tar.zst",
	"linux/amd64":   "cpython-3.12.5+20240814-x86_64-unknown-linux-gnu-pgo+lto-full.tar.zst",
	// "linux_gnu_x64_v2":  "cpython-3.12.5+20240814-x86_64_v2-unknown-linux-gnu-pgo+lto-full.tar.zst",
	// "linux_gnu_x64_v3":  "cpython-3.12.5+20240814-x86_64_v3-unknown-linux-gnu-pgo+lto-full.tar.zst",
	// "linux_gnu_x64_v4":  "cpython-3.12.5+20240814-x86_64_v4-unknown-linux-gnu-lto-full.tar.zst",
}

func (env *PyEnv) Install() {
	targetDir := filepath.Join(env.ParentPath, DIST_DIR)
	os.MkdirAll(targetDir, os.ModePerm)
	version := env.Distribution
	arch := versions[version]
	downloadPath := filepath.Join(targetDir, "python_download")
	downloadUrl := fmt.Sprintf("https://github.com/indygreg/python-build-standalone/releases/download/20240814/%s", arch)

	r, err := http.Get(downloadUrl)
	log.Printf("downloading embedded python tar from: %s\n", downloadUrl)
	if err != nil {
		log.Fatalf("download failed: %v\n", err)
	}
	if r.StatusCode == http.StatusNotFound {
		log.Fatalln("404 not found")
	}
	defer r.Body.Close()

	fileData, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("reading file data for write failed: %v\n", err)
	}
	log.Println("downloading embedded python tar complete")

	log.Printf("writing python tar to: %s\n", downloadPath)
	err = os.WriteFile(downloadPath, fileData, 0o640)
	if err != nil {
		err := os.RemoveAll(targetDir)
		if err != nil {
			log.Fatalf("removing bad download failed: %v\n", err)
		}
		log.Fatalf("writing file failed: %v\n", err)
	}
	log.Println("writing python tar complete")

	extract(downloadPath, targetDir)

	if strings.Contains(env.Distribution, "windows") {
		fp := filepath.Join(env.ParentPath, "dist/python/install/python.exe")
		log.Printf("installing pip to: %s\n", filepath.Join(env.ParentPath, "dist/python/install/Scripts"))
		err := installWindowsPip(fp)
		if err != nil {
			log.Fatalf("problem installing pip: %v\n", err)
		}
		log.Println("installing pip complete")
	}

	env.Compressed = false
	os.Remove(downloadPath)
}

func extract(archivePath string, targetPath string) string {
	log.Printf("extracting python tar to: %s\n", filepath.Join(targetPath, "python"))
	f, err := os.Open(archivePath)
	if err != nil {
		log.Fatalf("opening file failed: %v", err)
	}
	defer f.Close()

	z, err := zstd.NewReader(f)
	if err != nil {
		log.Fatalf("decompression failed: %v", err)
	}
	defer z.Close()
	// log.Printf("decompressing %s\n", archivePath)
	err = extractTarStream(z, targetPath)
	if err != nil {
		log.Fatalf("decompression failed: %v", err)
	}
	log.Println("extracting tar complete")

	return targetPath
}

func extractTarStream(r io.Reader, targetPath string) error {
	// log.Printf("extracting tar stream to: %s\n", targetPath)
	tarReader := tar.NewReader(r)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("ExtractTarStream: Next() failed: %w", err)
		}

		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q", header.Name)
		}

		p := filepath.FromSlash(header.Name)
		p = filepath.Join(targetPath, p)
		err = os.MkdirAll(filepath.Dir(p), 0755)
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(p, 0755); err != nil {
				return fmt.Errorf("ExtractTarStream: Mkdir() failed: %w", err)
			}
		case tar.TypeReg:
			_ = os.Remove(p) // we allow overwriting, which easily happens on case insensitive filesystems
			outFile, err := os.Create(p)
			if err != nil {
				return fmt.Errorf("ExtractTarStream: Create() failed: %w", err)
			}
			_, err = io.Copy(outFile, tarReader)
			_ = outFile.Close()
			if err != nil {
				return fmt.Errorf("ExtractTarStream: Copy() failed: %w", err)
			}
			err = os.Chmod(p, header.FileInfo().Mode())
			if err != nil {
				return fmt.Errorf("ExtractTarStream: Chmod() failed: %w", err)
			}
			err = os.Chtimes(p, header.AccessTime, header.ModTime)
			if err != nil {
				return err
			}
		case tar.TypeSymlink:
			_ = os.Remove(p) // we allow overwriting, which easily happens on case insensitive filesystems
			if err := os.Symlink(header.Linkname, p); err != nil {
				return fmt.Errorf("ExtractTarStream: Symlink() failed: %w", err)
			}
		default:
			return fmt.Errorf("ExtractTarStream: uknown type %v in %v", header.Typeflag, header.Name)
		}
	}
	return nil
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}

func installWindowsPip(fp string) error {
	// https://pip.pypa.io/en/stable/installation/
	cmd := exec.Command(fp, "-m", "ensurepip", "--upgrade")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
