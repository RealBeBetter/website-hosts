package core

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func AppendHostFile(filePath string) (err error) {
	// read content
	in, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file fail:", err)
		os.Exit(-1)
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(in)
	reader := bufio.NewReader(in)
	var strSlice []string
	line := 0
	startLine := 0
	endLine := 0
	for {
		line = line + 1
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		strSlice = append(strSlice, str)
		if str == startTag {
			startLine = line
		} else if str == endTag {
			endLine = line
		}
	}
	if startLine > 0 && endLine > 0 {
		strSlice = append(strSlice[:startLine-1], strSlice[endLine:]...)
	}
	str := strings.Join(strSlice, "")

	// write content
	err = CreateHostFile(str, filePath)
	return
}
