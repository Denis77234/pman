package fileHelper

import (
	"errors"
	"github.com/Masterminds/semver"
	"os"
	"path/filepath"
	"strings"
)

type Helper struct {
}

func (h *Helper) MakeDirIfNotExist(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *Helper) IsFileExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (h *Helper) Version(path string) string {

	str := strings.Replace(filepath.Base(path), ".zip", "", -1)

	return str
}

func (h *Helper) CheckVersion(file, lookingForVer string) (bool, error) {

	factVer := h.Version(file)

	con, err := semver.NewConstraint(lookingForVer)
	if err != nil {

		return false, err
	}

	v, err := semver.NewVersion(factVer)
	if err != nil {
		return false, err
	}

	b, err1 := con.Validate(v)
	if err != nil {
		return false, err1[0]
	}

	return b, nil

}
