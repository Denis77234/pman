package sshclient

import (
	"encoding/json"

	"fmt"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	reqstruct "packetManager/internal/Request"
	"packetManager/internal/fileHelper"
	"path/filepath"
	"strings"
	"time"
)

type Cfg struct {
	Username   string
	Password   string
	PrivateKey string
	Server     string
	Timeout    time.Duration
}

type Client struct {
	config     Cfg
	sshClient  *ssh.Client
	sftpClient *sftp.Client

	fileHelper.Helper
}

//---------------------------------------------------

func New(config Cfg) (*Client, error) {
	c := &Client{
		config: config,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) SendPack(sourcePath, destDir, packetName string) error {
	if err := c.connect(); err != nil {
		return fmt.Errorf("sendpack: %w", err)
	}

	dirPath := destDir + packetName

	err := c.sftpClient.Mkdir(dirPath)
	if err != nil {
		return err
	}

	filePath := dirPath + "/" + filepath.Base(sourcePath)

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

	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}

	return nil
}

func (c Client) DownloadPack(update reqstruct.Update, sourcePath, downloadTo string) error {
	if err := c.connect(); err != nil {
		return fmt.Errorf("downloadPack: %w", err)
	}

	for _, file := range update.Updates {
		path := sourcePath + "/" + file.Name
		fp, err := c.findZip(path, file.Ver)
		if err != nil {
			return err
		}

		err = c.downloadZip(fp, downloadTo)
		if err != nil {
			return err
		}

		if c.IsFileExist(path + "/" + "dependency.json") {

			jsonValue, err := os.Open(path + "/" + "dependency.json")
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

			err = c.DownloadPack(request, sourcePath, downloadTo)
			if err != nil {
				fmt.Println(err)
			}

		}
	}

	return nil
}

//-------------------------------------------------------------------

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

		valid, err := c.CheckVersion(file.Name(), ver)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "Invalid Semantic Version") {
				continue
			}
			return "", err
		}

		if valid {
			if validVersionFile == "" {
				validVersionFile = file.Name()
			} else {
				validVer := c.Version(validVersionFile)
				currentFileVer := c.Version(file.Name())
				compareStr := ">" + validVer

				bigger, _ := c.CheckVersion(currentFileVer, compareStr)
				if bigger {
					validVersionFile = file.Name()
				}

			}
		}
	}

	filePath := dir + "/" + validVersionFile
	return filePath, nil
}

func (c Client) downloadZip(file, downloadTo string) error {

	dirname := downloadTo + "/" + filepath.Base(filepath.Dir(file))

	c.MakeDirIfNotExist(dirname)

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
