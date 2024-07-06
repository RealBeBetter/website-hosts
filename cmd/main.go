package main

import "fmt"

import "github.com/RealBeBetter/website-hosts/internal/core"

func main() {
	fmt.Println("1. Create the host file to current directory (default)")
	fmt.Println("2. Append the host to system host file in linux (/etc/hosts)")
	fmt.Println("3. Append the host to system host file in windows (/mnt/c/Windows/System32/drivers/etc/hosts)")
	fmt.Println("4. Exit...")
	var mode int
	for {
		fmt.Print("Choose the mode [1-4]: ")
		_, _ = fmt.Scanln(&mode)
		if mode >= 0 && mode <= 4 {
			break
		}
	}

	var err error
	switch mode {
	case 0:
		fallthrough
	case 1:
		filePath := "./hosts"
		core.Backup(filePath)
		err = core.CreateHostFile("", filePath)
	case 2:
		filePath := "/etc/hosts"
		core.Backup(filePath)
		err = core.AppendHostFile(filePath)
	case 3:
		filePath := "/mnt/c/Windows/System32/drivers/etc/hosts"
		core.Backup(filePath)
		err = core.AppendHostFile(filePath)
	case 4:
		return
	}

	if err != nil {
		fmt.Println("error occurred:", err)
		return
	}

	fmt.Println("Done!")
}
