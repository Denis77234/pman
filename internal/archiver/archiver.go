package archiver

import (
	"archive/zip"
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
	packageDir string //path to directory where package to be archived is stored
	packageVer string
	archiveDir string // path to directory where the package should be archived
}

//-------------------------------------

func New(dir, ver, archDir string) Archiver {
	a := Archiver{packageDir: dir, packageVer: ver, archiveDir: "/home/denis/GolandProjects/packetManager/cmd/main"}
	return a
}

func (a Archiver) Archive() (archivePath string, err error) {
	dependencies, err := a.findDependencies()
	if err != nil {
		return "", err
	}

	archivePath = a.archiveDir + "/" + a.packageVer + ".zip"

	archive, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	defer zipWriter.Close()

	for _, dep := range dependencies {
		err = a.copyArchive(dep, zipWriter)
		if err != nil {
			return "", err
		}
	}

	err = a.copyArchive(a.pckg(), zipWriter)
	if err != nil {
		return "", err
	}

	return archivePath, nil
}

//-------------------------------------

func (a Archiver) pckg() string {
	pckg := a.packageDir + "/" + a.packageVer + ".zip"
	return pckg
}

func (a Archiver) dependencyPath(dependency string) string {

	depPath := fmt.Sprintf("%v/%v", filepath.Dir(a.packageDir), dependency)

	return depPath

}

func (a Archiver) dependencyFilePath() string {
	fileName := "dependency.json"
	depFilePath := fmt.Sprintf("%v/%v", a.packageDir, fileName)

	return depFilePath

}

func (a Archiver) isDirExist(path string) bool {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (a Archiver) version(path string) string {

	str := strings.Replace(filepath.Base(path), ".zip", "", -1)

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

	files, err := filepath.Glob(path + "/*.zip")
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

	return validVersionFile, nil
}

func (a Archiver) findDependencies() ([]string, error) {

	depFile, err := os.Open(a.dependencyFilePath())
	if err != nil {
		return nil, err
	}

	byteValue, err := io.ReadAll(depFile)
	if err != nil {
		return nil, err
	}
	dependency := make([]reqstruct.Packet, 0, 5)

	err = json.Unmarshal(byteValue, &dependency)
	if err != nil {
		return nil, err
	}

	depPackages := make([]string, 0, 10)

	for _, dep := range dependency {
		depPath := a.dependencyPath(dep.Name)

		if !a.isDirExist(depPath) {
			return nil, errors.New("dependency not found")
		}

		depPackPath, err := a.findDepPackage(depPath, dependency[0].Ver)
		if err != nil {
			return nil, err
		}

		depPackages = append(depPackages, depPackPath)

	}

	return depPackages, nil
}

func (a Archiver) copyArchive(from string, to *zip.Writer) error {

	archive, err := zip.OpenReader(from)
	if err != nil {
		return err
	}

	defer archive.Close()

	for _, file := range archive.File {

		packetName := filepath.Base(filepath.Dir(from))
		version := a.version(from)

		dir := packetName + "/" + version + "/"

		w, err := to.Create(dir + file.Name)
		if err != nil {
			return err
		}

		reader, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(w, reader)
		if err != nil {
			return err
		}
		err = reader.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
