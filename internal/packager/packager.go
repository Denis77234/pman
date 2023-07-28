package packager

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	reqstruct "packetManager/internal/Request"
)

type Packager struct {
	req         reqstruct.Request
	packagesDir string //path to directory where generated packages are stored
}

func New(request reqstruct.Request, packagesDir string) Packager {

	arch := Packager{req: request, packagesDir: packagesDir}

	return arch
}

// returns path for package directory
func (p Packager) pckDir() string {
	dir := fmt.Sprintf("%v/%v", p.packagesDir, p.req.Name)
	return dir
}

// returns path for archive file
func (p Packager) zipPath() string {
	zipPath := fmt.Sprintf("%v/%v", p.pckDir(), p.req.ArchiveName("tar"))
	return zipPath
}

func (p Packager) makeDirIfNotExist() error {
	if _, err := os.Stat(p.pckDir()); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(p.pckDir(), os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Packager) makeDependencyFile() error {

	fileName := "dependency.json"

	filePath := fmt.Sprintf("%v/%v", p.pckDir(), fileName)

	dep := p.req.Packets

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	enc := json.NewEncoder(file)

	enc.SetEscapeHTML(false)

	err = enc.Encode(dep)
	if err != nil {
		return err
	}

	return nil
}

func (p Packager) archiveMask(targ reqstruct.Target, zipWriter *zip.Writer) error {
	//search for files
	files, err := filepath.Glob(targ.Path)
	if err != nil {
		return err
	}

	for _, file := range files {
		// get file name
		fileName := filepath.Base(file)
		// if there is an exclusion filter...
		if targ.Exclude != "" {
			// ... check if the file name matches filter...
			exclude, err := filepath.Match(targ.Exclude, fileName)
			if err != nil {
				return err
			}
			//... and if it matches then don't handle it
			if exclude {
				continue
			}
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		w, err := zipWriter.Create(fileName)
		if err != nil {
			return err
		}

		_, err = io.Copy(w, f)
		if err != nil {
			return err
		}

		f.Close()
	}
	return nil
}

func (p Packager) Package() (packageDir, name string, err error) {

	err = p.makeDirIfNotExist()
	if err != nil {
		return "", "", err
	}

	arch, err := os.Create(p.zipPath())
	if err != nil {
		return "", "", err
	}

	defer arch.Close()

	zipWriter := zip.NewWriter(arch)

	defer zipWriter.Close()

	for _, target := range p.req.Targets {
		err := p.archiveMask(target, zipWriter)
		if err != nil {
			return "", "", err
		}
	}

	err = p.makeDependencyFile()
	if err != nil {
		return "", "", err
	}

	return p.pckDir(), p.req.Ver, nil
}
