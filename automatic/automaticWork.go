package automatic

import (
	"DocumentControl/service"
	"fmt"
	"github.com/dustin/go-humanize"
	"os"
	"strconv"
)

var item service.CopyFile

func Run() bool {
	return automaticOperation()
}

func automaticOperation() bool {
	fmt.Printf("=======================================\n")
	fmt.Printf("自动拷贝?\n")
	fmt.Printf("1. 是\n")
	fmt.Printf("2. 否\n")
	fmt.Printf("0. 退出\n")

	var input string
	fmt.Scanln(&input)
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("请输入序号\n")
		automaticOperation()
	}
	if num == 2 {
		return false
	} else if num == 0 {
		return true
	}

	automaticRun()
	automaticCopy()
	return true
}

func automaticRun() {
	item = service.CopyFile{}
	item.New()
	service.CheckDirectory(item.OptionData)
	item.SelectFiles = []service.Option{}
	item.RunReference()
	//createJSONFile()
}

func automaticCopy() {
	if len(item.SelectFiles) == 0 {
		return
	}
	for _, data := range item.SelectFiles {
		dataName, err := os.Stat(data.Address)
		if err != nil {
			fmt.Println(err)
		}

		item.FileTotal += dataName.Size()
	}
	fmt.Printf("\n开始复制:\n")
	service.FileTotal = item.FileTotal
	service.FileSize = 0

	item.CreateCopy(0)

	fmt.Println("\n完成,大小总共", humanize.Bytes(uint64(service.FileSize)))

	fmt.Printf("=======================================\n")
}
