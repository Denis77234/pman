package sshclient

import (
	"bytes"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestClient_DownloadPack(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			c.DownloadPack()
		})
	}
}

func TestClient_DownloadZip(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			if err := c.DownloadZip(); (err != nil) != tt.wantErr {
				t.Errorf("DownloadZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SendPack(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	type args struct {
		sourcePath string
		destDir    string
		packetName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			if err := c.SendPack(tt.args.sourcePath, tt.args.destDir, tt.args.packetName); (err != nil) != tt.wantErr {
				t.Errorf("SendPack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_checkVersion(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	type args struct {
		file          string
		lookingForVer string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			got, err := c.checkVersion(tt.args.file, tt.args.lookingForVer)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_close(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			c.close()
		})
	}
}

func TestClient_connect(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			if err := c.connect(); (err != nil) != tt.wantErr {
				t.Errorf("connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_create(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadWriteCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			got, err := c.create(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_download(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.ReadCloser
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			got, err := c.download(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("download() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_info(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    os.FileInfo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			got, err := c.info(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("info() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("info() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_upload(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	type args struct {
		source io.Reader
		size   int
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantDestination string
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			destination := &bytes.Buffer{}
			err := c.upload(tt.args.source, destination, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDestination := destination.String(); gotDestination != tt.wantDestination {
				t.Errorf("upload() gotDestination = %v, want %v", gotDestination, tt.wantDestination)
			}
		})
	}
}

func TestClient_version(t *testing.T) {
	type fields struct {
		config     Cfg
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	}
	type args struct {
		path string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				sshClient:  tt.fields.sshClient,
				sftpClient: tt.fields.sftpClient,
			}
			if got := c.version(tt.args.path); got != tt.want {
				t.Errorf("version() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		config Cfg
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
