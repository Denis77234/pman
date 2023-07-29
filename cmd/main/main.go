package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"packetManager/internal/Request"
	"packetManager/internal/archiver"
	"packetManager/internal/packager"
	"packetManager/internal/sshclient"
	"time"
)

func main() {

	jsonValue, err := os.Open("../../test/test.json")
	if err != nil {
		fmt.Println(err)
	}

	byteValue, err := io.ReadAll(jsonValue)
	if err != nil {
		fmt.Println(err)
	}

	request := reqstruct.Request{}

	jserr := json.Unmarshal(byteValue, &request)
	if jserr != nil {
		fmt.Println(jserr)
	}

	pack := packager.New(request, "/home/denis/GolandProjects/packetManager/packages")
	dir, ver, err := pack.Package()
	if err != nil {
		log.Fatal(err)
	}

	arch := archiver.New(dir, ver, "/home/denis/GolandProjects/packetManager/cmd/main/packet-1.zip")
	archivePath, err := arch.Archive()
	if err != nil {
		fmt.Println(err)
	}

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

	inf, err := cl.Info(archivePath)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(inf.Name())
}
