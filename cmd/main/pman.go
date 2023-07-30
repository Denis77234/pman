package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	reqstruct "packetManager/internal/Request"
	"packetManager/internal/archiver"
	"packetManager/internal/packager"
	"packetManager/internal/packetManager"
	"packetManager/internal/sshclient"
	"time"
)

const PACKAGESDIR = "/home/denis/GolandProjects/packetManager/packages"
const ARCHIVETO = "/home/denis/GolandProjects/packetManager/cmd"
const UPLOADTO = "/home/denis/dir/"

const DOWNLOADFROM = "/home/denis/sourceDir"
const DOWNLOADTO = "/home/denis/GolandProjects/packetManager/cmd/main"

func main() {

	cfgPM := packetManager.Config{
		DownloadFrom: DOWNLOADFROM,
		DownloadTo:   DOWNLOADTO,
		PackagesDir:  PACKAGESDIR,
		ArchiveTo:    ARCHIVETO,
		UploadTo:     UPLOADTO,
	}

	//jsonValue, err := os.Open("../../test/test.json")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//byteValue, err := io.ReadAll(jsonValue)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//request := reqstruct.Request{}
	//
	//jserr := json.Unmarshal(byteValue, &request)
	//if jserr != nil {
	//	fmt.Println(jserr)
	//}

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
		fmt.Println(err)
	}

	pm := packetManager.New(cfgPM, arch, pack, cl)

	jsonValue, err := os.Open("../../test/testdownload.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, err := io.ReadAll(jsonValue)
	if err != nil {
		fmt.Println(err)
	}

	request := reqstruct.Update{}

	jserr := json.Unmarshal(byteValue, &request)
	if jserr != nil {
		fmt.Println(jserr)
	}

	err = pm.DownloadPack(request)
	if err != nil {
		fmt.Println(err)
	}

}
