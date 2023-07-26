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

func (a Archiver) archiveMask(targ reqstruct.Target) error {
	files, err := filepath.Glob(targ.Path)
	if err != nil {
		return err
	}

	arch, err := os.Create(a.req.ArchiveName("tar"))
	if err != nil {
		fmt.Println(err)
	}
	defer arch.Close()

	zipWriter := zip.NewWriter(arch)

	defer zipWriter.Close()

	for _, file := range files {

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		w, err := zipWriter.Create(filepath.Base(file))
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

func (a Archiver) archiveExclude(targ reqstruct.Target) {

}

func (a Archiver) Archive() error {
	for _, target := range a.req.Targets {
		if target.Exclude != "" {
			continue
		}
		err := a.archiveMask(target)
		if err != nil {
			return err
		}
	}
	return nil
}
