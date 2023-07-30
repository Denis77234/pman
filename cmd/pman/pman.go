package main

import (
	"log"
	"packetManager/internal/archiver"
	"packetManager/internal/clidecorator"
	"packetManager/internal/packager"
	"packetManager/internal/packetManager"
	"packetManager/internal/sshclient"
	"time"
)

const PACKAGESDIR = "/home/denis/GolandProjects/packetManager/packages"
const ARCHIVETO = "/home/denis/GolandProjects/packetManager/cmd"
const UPLOADTO = "/home/denis/dir/"
const DOWNLOADFROM = "/home/denis/sourceDir"
const DOWNLOADTO = "/home/denis/GolandProjects/packetManager/cmd/pman"

func main() {

	cfgPM := packetManager.Config{
		DownloadFrom: DOWNLOADFROM,
		DownloadTo:   DOWNLOADTO,
		PackagesDir:  PACKAGESDIR,
		ArchiveTo:    ARCHIVETO,
		UploadTo:     UPLOADTO,
	}

	pack := packager.New()

	arch := archiver.New()

	cfg := sshclient.Cfg{
		Username: "denis",
		Password: "olodop73",
		Server:   "localhost:22",
		Timeout:  time.Second * 30,
	}

	cl, err := sshclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	pm := packetManager.New(cfgPM, arch, pack, cl)

	cli := clidecorator.CliDecorator{pm}

	cli.Start()

}
