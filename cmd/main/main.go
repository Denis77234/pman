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

	pack := packager.New(request, "../../packages")
	dir, ver, err := pack.Package()
	if err != nil {
		log.Fatal(err)
	}

	arch := archiver.New("../../packages", dir, ver)

	fmt.Println(arch.FindDependencies())

}
