package core

import (
	"fmt"
	"io"
	"os"
)

func Backup(filePath string) {
	filePathBak := filePath + "_bak"
	// backup
	backupErr := Copy(filePathBak, filePath)
	if backupErr != nil {
		fmt.Println("backup file fail:", backupErr)
		fmt.Println("continue...")
	}
}

func Copy(dstName, srcName string) (err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer func(src *os.File) {
		err = src.Close()
	}(src)
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func(dst *os.File) {
		err = dst.Close()
	}(dst)
	_, err = io.Copy(dst, src)
	return
}
