package main

import (
	"bufio"
	"github.com/riftbit/go-systray"
	"github.com/yob/graval"
	"io"
	"os"
	"time"
)

type PSFTPDriver struct{}

func (driver *PSFTPDriver) Authenticate(user string, pass string) bool {
	// Who?
	return user == User && pass == Pass
}

func (driver *PSFTPDriver) Bytes(path string) (bytes int64) {
	// Block Until Zip File Ready
	<-ZipFileReady

	// How?
	return ZipFileStat.Size()
}

func (driver *PSFTPDriver) ModifiedTime(path string) (time.Time, error) {
	// Block Until Zip File Ready
	<-ZipFileReady

	// When?
	return ZipFileStat.ModTime(), nil
}

func (driver *PSFTPDriver) ChangeDir(path string) bool {
	// Where?
	return path == "\\" || path == "/"
}

func (driver *PSFTPDriver) DirContents(path string) (files []os.FileInfo) {
	// Block Until Zip File Ready
	<-ZipFileReady

	// What?
	return append(files, graval.NewFileItem(Filename, ZipFileStat.Size(), ZipFileStat.ModTime()))
}

func (driver *PSFTPDriver) DeleteDir(path string) bool {
	return false
}

func (driver *PSFTPDriver) DeleteFile(path string) bool {
	return false
}

func (driver *PSFTPDriver) Rename(fromPath string, toPath string) bool {
	return false
}

func (driver *PSFTPDriver) MakeDir(path string) bool {
	return false
}

type PSFTPCloser struct {
	io.Reader
}

func (PSFTPCloser) Close() error {
	if *AutoQuit {
		_ = systray.Quit()
	}
	return nil
}

func (driver *PSFTPDriver) GetFile(path string) (reader io.ReadCloser, err error) {
	// Block Until Zip File is Ready
	<-ZipFileReady

	// Why Not Share the Zip!? :P
	zipFile, err := os.Open(ZipFile)
	if err != nil {
		reader = nil
	} else {
		reader = PSFTPCloser{bufio.NewReader(zipFile)}
	}
	return
}

func (driver *PSFTPDriver) PutFile(path string, data io.Reader) bool {
	return false
}

type PSFTPDriverFactory struct{}

func (factory *PSFTPDriverFactory) NewDriver() (graval.FTPDriver, error) {
	return &PSFTPDriver{}, nil
}
