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

func (a Archiver) version(path string) string {

	str := strings.Replace(filepath.Base(path), ".tar", "", -1)

	return str
}

func (a Archiver) checkVersion(file, lookingForVer string) (bool, error) {

	factVer := a.version(file)

	c, err := semver.NewConstraint(lookingForVer)
	if err != nil {

		return false, err
	}

	v, err := semver.NewVersion(factVer)
	if err != nil {
		return false, err
	}

	b, err1 := c.Validate(v)
	if err != nil {
		return false, err1[0]
	}

	return b, nil

}

func (a Archiver) findDepPackage(path, ver string) (string, error) {

	files, err := filepath.Glob(path + "/*.tar")
	if err != nil {
		return "", err
	}

	var validVersionFile string

	for _, file := range files {
		valid, err := a.checkVersion(file, ver)
		if err != nil {
			return "", err
		}
		if valid {
			if validVersionFile == "" {
				validVersionFile = file
			} else {
				validVer := a.version(validVersionFile)
				currentFileVer := a.version(file)
				compareStr := ">" + validVer

				bigger, _ := a.checkVersion(currentFileVer, compareStr)
				if bigger {
					validVersionFile = file
				}

			}
		}

	}

	fmt.Println(validVersionFile)

	return validVersionFile, nil
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

	depPackages := make([]string, 0, 10)

	for _, dep := range dependency {
		depPath := a.dependencyPath(dep.Name)

		if !a.isDependencyExist(depPath) {
			return errors.New("Dependecy not found")
		}

		depPackPath, err := a.findDepPackage(depPath, dependency[0].Ver)
		if err != nil {
			return err
		}

		depPackages = append(depPackages, depPackPath)

	}

	return nil
}
