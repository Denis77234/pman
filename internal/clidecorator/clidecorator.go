package clidecorator

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	reqstruct "packetManager/internal/Request"
	"packetManager/internal/packetManager"
)

type CliDecorator struct {
	packetManager.PacketManager
}

func (c CliDecorator) requestSend(path string) (reqstruct.Request, error) {
	jsonValue, err := os.Open(path)
	if err != nil {
		return reqstruct.Request{}, err
	}

	byteValue, err := io.ReadAll(jsonValue)
	if err != nil {
		return reqstruct.Request{}, err
	}

	request := reqstruct.Request{}

	jserr := json.Unmarshal(byteValue, &request)
	if jserr != nil {
		return reqstruct.Request{}, err
	}

	return request, nil
}

func (c CliDecorator) requestDownload(path string) (reqstruct.Update, error) {
	jsonValue, err := os.Open(path)
	if err != nil {
		return reqstruct.Update{}, err
	}

	byteValue, err := io.ReadAll(jsonValue)
	if err != nil {
		return reqstruct.Update{}, err
	}

	update := reqstruct.Update{}

	jserr := json.Unmarshal(byteValue, &update)
	if jserr != nil {
		return reqstruct.Update{}, err
	}

	return update, nil
}

func (c CliDecorator) Start() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  create  Create a package\n")
		fmt.Fprintf(os.Stderr, "  update  Update packages\n")
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Arg(0)

	path := flag.Arg(1)

	switch command {
	case "create":

		if flag.NArg() < 2 {
			fmt.Fprintf(os.Stderr, "Usage: %s create 'path to json file'\n", os.Args[0])
			os.Exit(1)
		}

		requestSend, err := c.requestSend(path)
		if err != nil {
			log.Fatal(err)
		}

		err = c.MakeAndSendPackage(requestSend)
		if err != nil {
			log.Fatal(err)
		}
	case "update":

		if flag.NArg() < 2 {
			fmt.Fprintf(os.Stderr, "Usage: %s create 'path to json file'\n", os.Args[0])
			os.Exit(1)
		}

		downloadReq, err := c.requestDownload(path)

		err = c.DownloadPack(downloadReq)
		if err != nil {
			log.Fatal(err)
		}
	default:
		fmt.Fprintln(os.Stderr, "Unknown command. Available commands:")
		flag.Usage()
		os.Exit(1)
	}
}
