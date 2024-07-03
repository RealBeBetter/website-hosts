package core

import (
	"bufio"
	"fmt"
	"os"
)

func CreateHostFile(str string, filePath string) (err error) {
	str += startTag
	ch := make(chan *HostChan)
	for _, v := range domains {
		go GetWebsiteHost(v, ch)
	}
	fmt.Println("================\nstart get host：\n================")
	hostMap := make(map[string]string)
	for range domains {
		chRec := <-ch
		if chRec.Err != nil {
			fmt.Println(chRec.Err.Error() + " " + chRec.Domain)
			err = chRec.Err
			return
		}
		hostMap[chRec.Domain] = chRec.Ip
		fmt.Println(chRec.Ip + " " + chRec.Domain)
	}
	for _, v := range domains {
		str += hostMap[v] + " " + v + "\r\n"
	}

	str += endTag
	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("open file fail:", err)
		return
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)
	writer := bufio.NewWriter(out)
	_, _ = writer.WriteString(str)
	_ = writer.Flush()
	fmt.Println("================\ndone\n================")
	return
}
