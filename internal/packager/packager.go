package packager

import (
	"archive/zip"
	"encoding/json"
	"io"
	"os"
	"packetManager/internal/fileHelper"
	"path/filepath"

	reqstruct "packetManager/internal/Request"
)

type Packager struct {
	//req         reqstruct.Request
	//packagesDir string //path to directory where generated packages are stored

	fileHelper.Helper
}

//----------------------------------------------------

func New() Packager {

	arch := Packager{}

	return arch
}

func (p Packager) Package(req reqstruct.Request, packagesDir string) (pckDir, pckVer string, err error) {

	pckDir = filepath.Join(packagesDir, req.Name)

	err = p.MakeDirIfNotExist(pckDir)
	if err != nil {
		return "", "", err
	}

	arch, err := os.Create(p.zipPath(req, packagesDir))
	if err != nil {
		return "", "", err
	}

	defer arch.Close()

	zipWriter := zip.NewWriter(arch)

	defer zipWriter.Close()

	for _, target := range req.Targets {
		err := p.archiveMask(target, zipWriter)
		if err != nil {
			return "", "", err
		}
	}

	err = p.makeDependencyFile(req, packagesDir)
	if err != nil {
		return "", "", err
	}

	return pckDir, req.Ver, nil
}

//------------------------------------------------------

// returns path for archive file
func (p Packager) zipPath(req reqstruct.Request, packagesDir string) string {

	zipPath := filepath.Join(packagesDir, req.Name, req.ArchiveName("zip"))

	return zipPath
}

func (p Packager) makeDependencyFile(req reqstruct.Request, packagesDir string) error {

	fileName := "dependency.json"

	filePath := filepath.Join(packagesDir, req.Name, fileName)

	dep := req.Packets

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
