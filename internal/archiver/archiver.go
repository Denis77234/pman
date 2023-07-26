package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	reqstruct "packetManager/internal/Request"
	"path/filepath"
)

type Archiver struct {
	req reqstruct.Request
}

func New(request reqstruct.Request) Archiver {

	arch := Archiver{req: request}

	return arch
}

func (a Archiver) archiveMask(targ reqstruct.Target, zipWriter *zip.Writer) error {
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

func (a Archiver) Archive() error {

	arch, err := os.Create(a.req.ArchiveName("tar"))
	if err != nil {
		fmt.Println(err)
	}

	defer arch.Close()

	zipWriter := zip.NewWriter(arch)

	defer zipWriter.Close()

	for _, target := range a.req.Targets {
		err := a.archiveMask(target, zipWriter)
		if err != nil {
			return err
		}
	}
	return nil
}
