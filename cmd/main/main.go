package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	reqstruct "packetManager/internal/Request"
	"packetManager/internal/sshclient"
	"time"
)

func main() {

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
	//
	//pack := packager.New(request, "/home/denis/GolandProjects/packetManager/packages")
	//dir, err := pack.Package()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//arch := archiver.New(dir, request.Ver, "/home/denis/GolandProjects/packetManager/cmd")
	//_, err = arch.Archive()
	//if err != nil {
	//	fmt.Println(err)
	//}

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

	//sendPath := "/home/denis/dir/"

	//err = cl.SendPack(archivePath, sendPath, request.Name)
	//if err != nil {
	//	fmt.Println(err)
	//}

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

	source := "/home/denis/sourceDir"
	err = cl.DownloadPack(request, source)
	if err != nil {
		fmt.Println(err)
	}
}
