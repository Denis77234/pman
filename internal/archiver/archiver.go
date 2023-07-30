package archiver

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	reqstruct "packetManager/internal/Request"
	"packetManager/internal/fileHelper"
	"path/filepath"
)

type Archiver struct {
	fileHelper.Helper
}

//-------------------------------------

func New() Archiver {
	a := Archiver{}
	return a
}

func (a Archiver) Archive(packageDir, packageVer, archiveTo string) (archivePath string, err error) {
	dependencies, err := a.findDependencies(packageDir)
	if err != nil {
		return "", err
	}

	archivePath = archiveTo + "/" + packageVer + ".zip"

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

	err = a.copyArchive(a.pckg(packageDir, packageVer), zipWriter)
	if err != nil {
		return "", err
	}

	return archivePath, nil
}

//-------------------------------------

func (a Archiver) pckg(packageDir, packageVer string) string {
	pckg := packageDir + "/" + packageVer + ".zip"
	return pckg
}

func (a Archiver) dependencyPath(dependency string, packageDir string) string {

	depPath := fmt.Sprintf("%v/%v", filepath.Dir(packageDir), dependency)

	return depPath

}

func (a Archiver) dependencyFilePath(packageDir string) string {
	fileName := "dependency.json"
	depFilePath := fmt.Sprintf("%v/%v", packageDir, fileName)

	return depFilePath

}

func (a Archiver) findDepPackage(path, ver string) (string, error) {

	files, err := filepath.Glob(path + "/*.zip")
	if err != nil {
		return "", err
	}

	var validVersionFile string

	for _, file := range files {
		valid, err := a.CheckVersion(file, ver)
		if err != nil {
			return "", err
		}
		if valid {
			if validVersionFile == "" {
				validVersionFile = file
			} else {
				validVer := a.Version(validVersionFile)
				currentFileVer := a.Version(file)
				compareStr := ">" + validVer

				bigger, _ := a.CheckVersion(currentFileVer, compareStr)
				if bigger {
					validVersionFile = file
				}

			}
		}

	}

	return validVersionFile, nil
}

func (a Archiver) findDependencies(packageDir string) ([]string, error) {

	depFile, err := os.Open(a.dependencyFilePath(packageDir))
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
		depPath := a.dependencyPath(dep.Name, packageDir)

		if !a.IsFileExist(depPath) {
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
		version := a.Version(from)

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
