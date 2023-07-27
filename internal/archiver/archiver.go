package archiver

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"io"
	"os"
	"path/filepath"
	"strings"

	reqstruct "packetManager/internal/Request"
)

type Archiver struct {
	packagesDir string //path to directory where generated packages are stored
	packageDir  string //path to directory where package to be archived is stored
	packageVer  string
}

func New(pDir, dir, ver string) Archiver {
	a := Archiver{packagesDir: pDir, packageDir: dir, packageVer: ver}
	return a
}

func (a Archiver) dependencyPath(dependency string) string {

	depPath := fmt.Sprintf("%v/%v", a.packagesDir, dependency)

	return depPath

}

func (a Archiver) dependencyFilePath() string {
	fileName := "dependency.json"
	depFilePath := fmt.Sprintf("%v/%v", a.packageDir, fileName)

	return depFilePath

}

func (a Archiver) isDependencyExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (a Archiver) version(path string) (*semver.Version, error) {

	str := strings.Replace(filepath.Base(path), ".tar", "", -1)

	ver, err := semver.NewVersion(str)
	if err != nil {
		return nil, err
	}
	return ver, nil
}

//func (a Archiver) checkVerson() {
//
//}

func (a Archiver) findDepPackage(path, ver string) error {

	files, err := filepath.Glob(path + "/*.tar")
	if err != nil {
		return err
	}

	fmt.Println(a.version(files[0]))

	return nil
}

func (a Archiver) FindDependencies() error {

	depFile, err := os.Open(a.dependencyFilePath())
	if err != nil {
		return err
	}

	byteValue, err := io.ReadAll(depFile)
	if err != nil {
		return err
	}
	dependency := []reqstruct.Packet{}

	err = json.Unmarshal(byteValue, &dependency)
	if err != nil {
		return err
	}

	for _, dep := range dependency {
		depPath := a.dependencyPath(dep.Name)

		if !a.isDependencyExist(depPath) {
			return errors.New("Dependecy not found")
		}

		a.findDepPackage(depPath, a.packageVer)
		
	}

	return nil
}
