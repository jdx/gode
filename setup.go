package gode

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// IsSetup returns true if node is setup in the client's RootPath directory
func (c *Client) IsSetup() bool {
	// TODO: better check if it is setup
	exists, _ := fileExists(c.nodePath())
	return exists
}

// Setup downloads and sets up node in the client's RootPath directory
func (c *Client) Setup() error {
	if runtime.GOOS == "windows" {
		return c.setupWindows()
	}
	return c.setupUnix()
}

func (c *Client) setupUnix() error {
	err := os.MkdirAll(c.RootPath, 0777)
	if err != nil {
		return err
	}
	resp, err := http.Get(c.nodeURL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	uncompressed, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	return extractArchive(tar.NewReader(uncompressed), c.RootPath)
}

func extractArchive(archive *tar.Reader, rootPath string) error {
	for {
		hdr, err := archive.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		path := filepath.Join(rootPath, hdr.Name)
		switch {
		case hdr.FileInfo().IsDir():
			if err := os.Mkdir(path, hdr.FileInfo().Mode()); err != nil {
				return err
			}
		case hdr.Linkname != "":
			if err := os.Symlink(hdr.Linkname, path); err != nil {
				return err
			}
		default:
			if err := extractFile(archive, hdr, path); err != nil {
				return err
			}
		}
	}
}

func extractFile(archive *tar.Reader, hdr *tar.Header, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, archive)
	if err != nil {
		return err
	}
	return os.Chmod(path, hdr.FileInfo().Mode())
}

func (c *Client) setupWindows() error {
	err := c.downloadNodeExe()
	if err != nil {
		return err
	}
	return c.downloadNpm()
}

func (c *Client) downloadNodeExe() error {
	err := os.MkdirAll(filepath.Dir(c.nodePath()), 0777)
	if err != nil {
		return err
	}
	resp, err := http.Get(c.nodeURL())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(c.nodePath())
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	return err
}

func (c *Client) downloadNpm() error {
	modulesDir := filepath.Join(c.RootPath, c.nodeBase(), "lib", "node_modules")
	err := os.MkdirAll(modulesDir, 0777)
	if err != nil {
		return err
	}
	resp, err := http.Get("https://github.com/npm/npm/archive/v2.1.6.tar.gz")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	uncompressed, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	err = extractArchive(tar.NewReader(uncompressed), modulesDir)
	if err != nil {
		return err
	}
	return os.Rename(filepath.Join(modulesDir, "npm-2.1.6"), filepath.Join(modulesDir, "npm"))
}
