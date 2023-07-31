package envhelper

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"packetManager/internal/packetManager"
	"packetManager/internal/sshclient"
	"time"
)

type EnvHelper struct {
}

func (e EnvHelper) LoadEnvForDirs(envFilePath string) error {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return err
	}

	err = e.validDirEnv()
	if err != nil {
		return err
	}

	return nil
}

func (e EnvHelper) MakePMConfig() (cfg packetManager.Config) {
	cfg.DownloadTo = os.Getenv("PMAN_DOWNLOADTO")
	cfg.ArchiveTo = os.Getenv("PMAN_ARCHIVETO")
	cfg.DownloadFrom = os.Getenv("PMAN_DOWNLOADFROM")
	cfg.UploadTo = os.Getenv("PMAN_UPLOADTO")
	cfg.PackagesDir = os.Getenv("PMAN_PACKAGESDIR")
	return cfg
}

func (e EnvHelper) validDirEnv() error {
	pd := os.Getenv("PMAN_PACKAGESDIR")
	arc := os.Getenv("PMAN_ARCHIVETO")
	up := os.Getenv("PMAN_UPLOADTO")
	df := os.Getenv("PMAN_DOWNLOADFROM")
	dt := os.Getenv("PMAN_DOWNLOADTO")
	if pd == "" || arc == "" || up == "" || df == "" || dt == "" {
		return errors.New("invalid dir env file")
	}

	return nil
}

func (e EnvHelper) LoadEnvForSSH(envFilePath string) error {
	err := godotenv.Load(envFilePath)
	if err != nil {
		return err
	}

	err = e.validSSHEnv()
	if err != nil {
		return err
	}

	return nil
}

func (e EnvHelper) MakeSSHConfig() (cfg sshclient.Cfg) {
	cfg.Username = os.Getenv("PMAN_USERNAME")
	cfg.Password = os.Getenv("PMAN_PASSSWORD")
	cfg.PrivateKey = os.Getenv("PMAN_PRIVATE_KEY")
	cfg.Server = os.Getenv("PMAN_SERVER")
	cfg.Timeout = time.Second * 30
	return cfg
}

func (e EnvHelper) validSSHEnv() error {
	usr := os.Getenv("PMAN_USERNAME")
	pw := os.Getenv("PMAN_PASSSWORD")
	pk := os.Getenv("PMAN_PRIVATE_KEY")
	srv := os.Getenv("PMAN_SERVER")
	if usr == "" || (pw == "" && pk == "") || srv == "" {
		return errors.New("invalid ssh env file")
	}

	return nil
}
