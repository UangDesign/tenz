package tenz

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/UangDesign/filetype"
)

type ITenZ interface {
	// 解压缩
	DeCompress(scrFile string, dirName string) error
	// 压缩
	Compress() error
}

type TenZ struct {
}

func NewTenZ() ITenZ {
	return &TenZ{}
}

const (
	FILE_TYPE_ZIP    = ".zip"
	FILE_TYPE_TAR    = ".tar"
	FILE_TYPE_TAR_GZ = ".tar.gz"
	FILE_TYPE_TGZ    = ".tgz"
	FILE_TYPE_GZ     = ".gz"
)

/*
	@srcFile	string	原文件
	@dirname	string	目标地址
	@fileType	string	文件类型
*/
func (s *TenZ) DeCompress(srcFile, dirname string) (err error) {
	switch filetype.NewFileType().GetFileType(srcFile) {
	case FILE_TYPE_ZIP:
		err = unPackZIP(srcFile, dirname)
	case FILE_TYPE_TAR:
		err = unPackTAR(srcFile, dirname)
	case FILE_TYPE_TAR_GZ:
	case FILE_TYPE_TGZ:
		err = unPackTGZ(srcFile, dirname)
	case FILE_TYPE_GZ:
		err = unPackGZ(srcFile, dirname)
	}
	return err
}
func (s *TenZ) Compress() (err error) {
	return fmt.Errorf("Can not compress")
}

func unPackZIP(srcFile, dirname string) (err error) {
	reader, err := zip.OpenReader(srcFile)
	defer reader.Close()
	if err != nil {
		log.Fatalf("open zip file to reader failed, err is: %v", err)
	} else {
		for _, zipfile := range reader.File {
			if !zipfile.FileInfo().IsDir() {
				rc, err := zipfile.Open()
				defer rc.Close()
				if err != nil {
					err = fmt.Errorf("Open file failed, err is: %v", err)
				} else {
					file, err := os.Create(path.Join(dirname, zipfile.Name))
					defer file.Close()
					if err != nil {
						err = fmt.Errorf("unzip to create file failed, err is: %c", err)
					} else {
						if _, err := io.Copy(file, rc); err != nil {
							err = fmt.Errorf("unzip to copy file failed, err is:%v", err)
						}
					}
				}
			} else {
				if err := os.MkdirAll(path.Join(dirname, zipfile.Name), os.FileMode(0755)); err != nil {
					err = fmt.Errorf("unpack zip to create dir failed, err is: %v", err)
				}
			}
		}
	}
	return err
}
func unPackTAR(srcFile, dirname string) (err error) {
	tarFile, err := os.Open(srcFile)
	defer tarFile.Close()
	if err != nil {
		err = fmt.Errorf("Open tar file failed, err is: %v", err)
	} else {
		err = tarReaderFunc(tar.NewReader(tarFile), dirname)
	}
	return err
}

func unPackTGZ(srcFile, dirname string) (err error) {
	if tgzFile, err := os.Open(srcFile); err != nil {
		err = fmt.Errorf("Open tar file failed, err is: %v", err)
		defer tgzFile.Close()
	} else {
		tgzReader, err := gzip.NewReader(tgzFile)
		defer tgzReader.Close()
		if err != nil {
			err = fmt.Errorf("get tgz reader filed, err is:%v", err)
		} else {
			err = tarReaderFunc(tar.NewReader(tgzReader), dirname)
		}
	}
	return err
}

func unPackGZ(srcFile, dirname string) (err error) {
	gzFile, err := os.Open(srcFile)
	defer gzFile.Close()
	if err != nil {
		err = fmt.Errorf("Open tar file failed, err is: %v", err)
	} else {
		gzReader, err := gzip.NewReader(gzFile)
		defer gzReader.Close()
		if err != nil {
			err = fmt.Errorf("new gz reader fialed, err is:%v", err)
		} else {
			filename := path.Join(dirname, gzReader.Name)
			if file, err := os.Create(filename); err != nil {
				err = fmt.Errorf("create file to copy failed, err: %v", err)
			} else {
				if _, err := io.Copy(file, gzReader); err != nil {
					err = fmt.Errorf("unzip to copy file failed, err is:%v", err)
				}
			}
		}
	}
	return err
}

func tarReaderFunc(tarReader *tar.Reader, dirname string) (err error) {
	for {
		tarHeader, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				err = fmt.Errorf("tarReader next failed, err is:%v", err)
			}
		}
		filename := path.Join(dirname, tarHeader.Name)
		if tarHeader.FileInfo().IsDir() {
			if err := os.MkdirAll(filename, os.FileMode(0755)); err != nil {
				err = fmt.Errorf("unpack tar to crate dir failed, err is:%v", err)
			}
		} else {
			pathDirName := path.Dir(filename)
			if _, err := os.Lstat(pathDirName); err != nil {
				if err := os.MkdirAll(pathDirName, os.FileMode(0755)); err != nil {
					err = fmt.Errorf("unpack tar to crate dir failed, err is:%v", err)
				}
			}
			file, err := os.Create(filename)
			defer file.Close()
			if err != nil {
				err = fmt.Errorf("create file to copy failed, err is:%v", err)
			} else {
				if _, err := io.Copy(file, tarReader); err != nil {
					err = fmt.Errorf("unzip to copy file failed, err is:%v", err)
				}
			}
		}
	}
	return err
}
