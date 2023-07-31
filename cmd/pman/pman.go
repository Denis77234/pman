package main

import (
	"log"
	"packetManager/internal/archiver"
	"packetManager/internal/clidecorator"
	"packetManager/internal/envhelper"
	"packetManager/internal/packager"
	"packetManager/internal/packetManager"
	"packetManager/internal/sshclient"
)

func main() {

	ev := envhelper.EnvHelper{}

	err := ev.LoadEnvForDirs("directories.env")
	if err != nil {
		log.Fatal(err)
	}

	err = ev.LoadEnvForSSH("sshCFG.env")
	if err != nil {
		log.Fatal(err)
	}

	cfgPM := ev.MakePMConfig()

	pack := packager.New()

	arch := archiver.New()

	cfg := ev.MakeSSHConfig()

	cl, err := sshclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	pm := packetManager.New(cfgPM, arch, pack, cl)

	cli := clidecorator.CliDecorator{PacketManager: pm}

	cli.Start()

}
