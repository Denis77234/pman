package packetManager

import (
	reqstruct "packetManager/internal/Request"
)

type abstractArchiver interface {
	Archive(packageDir, packageVer, archiveTo string) (archivePath string, err error)
}

type abstractPackager interface {
	Package(req reqstruct.Request, packagesDir string) (pckDir, pckVer string, err error)
}

type senderDownloader interface {
	SendPack(sourcePath, destDir, packetName string) error
	DownloadPack(update reqstruct.Update, sourcePath, downloadTo string) error
}

type Config struct {
	DownloadFrom string
	DownloadTo   string
	PackagesDir  string //path to directory where generated packages are stored
	ArchiveTo    string
	UploadTo     string
}

type PacketManager struct {
	cfg      Config
	archiver abstractArchiver
	packager abstractPackager
	client   senderDownloader
}

func New(cfg Config, archiver abstractArchiver, packager abstractPackager, client senderDownloader) PacketManager {

	return PacketManager{cfg: cfg, archiver: archiver, packager: packager, client: client}

}

func (p PacketManager) MakeAndSendPackage(requestSend reqstruct.Request) error {
	dir, ver, err := p.packager.Package(requestSend, p.cfg.PackagesDir)
	if err != nil {
		return err
	}
	archPath, err := p.archiver.Archive(dir, ver, p.cfg.ArchiveTo)
	if err != nil {
		return err
	}

	err = p.client.SendPack(archPath, p.cfg.UploadTo, requestSend.Name)
	if err != nil {
		return err
	}

	return nil
}

func (p PacketManager) DownloadPack(downloadReq reqstruct.Update) error {

	err := p.client.DownloadPack(downloadReq, p.cfg.DownloadFrom, p.cfg.DownloadTo)
	if err != nil {
		return err
	}
	return err
}
