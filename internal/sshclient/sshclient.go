package sshclient

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	reqstruct "packetManager/internal/Request"
	"path/filepath"
	"strings"
	"time"
)

type Cfg struct {
	Username     string
	Password     string
	PrivateKey   string
	Server       string
	KeyExchanges []string
	Timeout      time.Duration
}

type Client struct {
	config Cfg

	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

func New(config Cfg) (*Client, error) {
	c := &Client{
		config: config,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) makeDirIfNotExist(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) version(path string) string {

	str := strings.Replace(filepath.Base(path), ".zip", "", -1)

	return str
}

func (c *Client) checkVersion(file, lookingForVer string) (bool, error) {

	factVer := c.version(file)

	con, err := semver.NewConstraint(lookingForVer)
	if err != nil {

		return false, err
	}

	v, err := semver.NewVersion(factVer)
	if err != nil {
		return false, err
	}

	b, err1 := con.Validate(v)
	if err != nil {
		return false, err1[0]
	}

	return b, nil

}

func (c *Client) info(filePath string) (os.FileInfo, error) {
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	info, err := c.sftpClient.Lstat(filePath)
	if err != nil {
		return nil, fmt.Errorf("file stats: %w", err)
	}

	return info, nil
}

func (c *Client) close() {
	if c.sftpClient != nil {
		c.sftpClient.Close()
	}
	if c.sshClient != nil {
		c.sshClient.Close()
	}
}

func (c *Client) connect() error {
	if c.sshClient != nil {
		_, _, err := c.sshClient.SendRequest("keepalive", false, nil)
		if err == nil {
			return nil
		}
	}

	auth := ssh.Password(c.config.Password)
	if c.config.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(c.config.PrivateKey))
		if err != nil {
			return fmt.Errorf("ssh parse private key: %w", err)
		}
		auth = ssh.PublicKeys(signer)
	}

	cfg := &ssh.ClientConfig{
		User: c.config.Username,
		Auth: []ssh.AuthMethod{
			auth,
		},
		HostKeyCallback: func(string, net.Addr, ssh.PublicKey) error { return nil },
		Timeout:         c.config.Timeout,
		Config: ssh.Config{
			KeyExchanges: c.config.KeyExchanges,
		},
	}

	sshClient, err := ssh.Dial("tcp", c.config.Server, cfg)
	if err != nil {
		return fmt.Errorf("ssh dial: %w", err)
	}
	c.sshClient = sshClient

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("sftp new client: %w", err)
	}
	c.sftpClient = sftpClient

	return nil
}

func (c *Client) upload(source io.Reader, destination io.Writer, size int) error {
	if err := c.connect(); err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	chunk := make([]byte, size)

	for {
		num, err := source.Read(chunk)
		if err == io.EOF {
			tot, err := destination.Write(chunk[:num])
			if err != nil {
				return err
			}

			if tot != len(chunk[:num]) {
				return fmt.Errorf("failed to write stream")
			}

			return nil
		}

		if err != nil {
			return err
		}

		tot, err := destination.Write(chunk[:num])
		if err != nil {
			return err
		}

		if tot != len(chunk[:num]) {
			return fmt.Errorf("failed to write stream")
		}
	}
}

func (c *Client) create(filePath string) (io.ReadWriteCloser, error) {
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return c.sftpClient.Create(filePath)
}

func (c *Client) SendPack(sourcePath, destDir, packetName string) error {

	dirPath := destDir + packetName

	err := c.sftpClient.Mkdir(dirPath)
	if err != nil {
		return err
	}

	filePath := dirPath + "/" + filepath.Base(sourcePath)
	fmt.Println(filePath)

	dest, err := c.create(filePath)
	if err != nil {
		return err
	}

	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	defer source.Close()

	err = c.upload(source, dest, 10)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) download(filePath string) (io.ReadCloser, error) {
	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	return c.sftpClient.Open(filePath)
}

func (c Client) findZip(dir, ver string) (string, error) {

	if ver == "" {
		ver = ">=0.0.1"
	}

	files, err := c.sftpClient.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var validVersionFile string

	for _, file := range files {
		valid, err := c.checkVersion(file.Name(), ver)
		if err != nil {
			return "", err
		}

		if valid {
			if validVersionFile == "" {
				validVersionFile = file.Name()
			} else {
				validVer := c.version(validVersionFile)
				currentFileVer := c.version(file.Name())
				compareStr := ">" + validVer

				bigger, _ := c.checkVersion(currentFileVer, compareStr)
				if bigger {
					validVersionFile = file.Name()
				}

			}
		}
	}

	filePath := dir + "/" + validVersionFile
	return filePath, nil
}

func (c Client) downloadZip(file string) error {

	dirname := filepath.Base(filepath.Dir(file))

	c.makeDirIfNotExist(dirname)

	f, err := os.Create(dirname + "/" + filepath.Base(file))
	if err != nil {
		return err
	}

	d, err := c.download(file)
	if err != nil {

		return err
	}

	defer d.Close()

	bytes, err := io.ReadAll(d)
	if err != nil {

		return err
	}

	_, err = f.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) DownloadPack(update reqstruct.Update, soursePath string) error {

	for _, file := range update.Updates {
		path := soursePath + "/" + file.Name
		fp, err := c.findZip(path, file.Version)
		if err != nil {
			return err
		}

		c.downloadZip(fp)
	}
	return nil
}
