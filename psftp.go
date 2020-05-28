package main

import (
	"bufio"
	"github.com/riftbit/go-systray"
	"github.com/yob/graval"
	"io"
	"os"
	"time"
)

type PSFTPDriver struct {
	User         string
	Pass         string
	Filename     string
	ZipFileReady chan bool
	ZipFileStat  os.FileInfo
	ZipFile      os.File
}

func (driver *PSFTPDriver) Authenticate(user string, pass string) bool {
	// Who?
	return user == driver.User && pass == driver.Pass
}

func (driver *PSFTPDriver) Bytes(path string) (bytes int64) {
	// Block Until Zip File Ready
	<-driver.ZipFileReady

	// How?
	return driver.ZipFileStat.Size()
}

func (driver *PSFTPDriver) ModifiedTime(path string) (time.Time, error) {
	// Block Until Zip File Ready
	<-driver.ZipFileReady

	// When?
	return driver.ZipFileStat.ModTime(), nil
}

func (driver *PSFTPDriver) ChangeDir(path string) bool {
	// Where?
	return path == "\\" || path == "/"
}

func (driver *PSFTPDriver) DirContents(path string) (files []os.FileInfo) {
	// Block Until Zip File Ready
	<-driver.ZipFileReady

	// What?
	return append(files, graval.NewFileItem(driver.Filename, driver.ZipFileStat.Size(), driver.ZipFileStat.ModTime()))
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
	AutoQuit *bool
}

func (closer PSFTPCloser) Close() error {
	if *closer.AutoQuit {
		_ = systray.Quit()
	}
	return nil
}

func (driver *PSFTPDriver) GetFile(path string) (reader io.ReadCloser, err error) {
	// Block Until Zip File is Ready
	<-driver.ZipFileReady

	// Why Not Share the Zip!? :P
	zipFile, err := os.Open(driver.ZipFile)
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
