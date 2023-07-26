package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"packetManager/internal/Request"
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

	arch := packager.New(request, "../../packages")
	err = arch.Archive()
	if err != nil {
		log.Fatal(err)
	}
}
