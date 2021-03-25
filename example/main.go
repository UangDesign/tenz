package main

import (
	"log"
	"path/filepath"

	"github.com/UangDesign/tenz"
)

var (
	// ZIP
	FILE_SRC_PATH_ZIP = getCurrentPath("./file/zip")
	ZIP_PATH          = filepath.Join(FILE_SRC_PATH_ZIP, "test.zip")
	// TAR
	FILE_SRC_PATH_TAR = getCurrentPath("./file/tar")
	TAR_PATH          = filepath.Join(FILE_SRC_PATH_TAR, "test.tar")
	// TGZ
	FILE_SRC_PATH_TGZ = getCurrentPath("./file/tgz")
	TGZ_PATH          = filepath.Join(FILE_SRC_PATH_TGZ, "test.tgz")
	// GZ
	FILE_SRC_PATH_GZ = getCurrentPath("./file/gz")
	GZ_PATH          = filepath.Join(FILE_SRC_PATH_GZ, "test.gz")
)

func main() {
	var err error
	deCompressObj := tenz.NewTenZ()
	// zip
	err = deCompressObj.DeCompress(ZIP_PATH, FILE_SRC_PATH_ZIP)
	if err != nil {
		log.Fatalf("decompress zip file failed, err: %v", err)
	}
	// tgz
	err = deCompressObj.DeCompress(TGZ_PATH, FILE_SRC_PATH_TGZ)
	if err != nil {
		log.Fatalf("decompress tgz file failed, err: %v", err)
	}
	// gz
	err = deCompressObj.DeCompress(GZ_PATH, FILE_SRC_PATH_GZ)
	if err != nil {
		log.Fatalf("decompress gz file failed, err: %v", err)
	}
	// tar
	err = deCompressObj.DeCompress(TAR_PATH, FILE_SRC_PATH_TAR)
	if err != nil {
		log.Fatalf("decompress tar file failed, err: %v", err)
	}
}

func getCurrentPath(dir string) (absPath string) {
	if curPath, err := filepath.Abs(dir); err == nil {
		absPath = curPath
	}
	return absPath
}
