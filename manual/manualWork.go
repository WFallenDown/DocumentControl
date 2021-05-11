package manual

import (
	"DocumentControl/service"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type FilesData struct {
	FileName string
	FileType string
	Tag      []string
	Describe string
	Address  string
	Child    []FilesData
}

var number int
var address string
var item service.CopyFile

func SelectFolder() {
	fmt.Printf("=======================================\n")
	fmt.Printf("请选择需要对比的文件夹(输入序号):\n")

	fmt.Printf("1. Animations\n")
	fmt.Printf("2. Flicks\n")
	fmt.Printf("3. CG\n")
	fmt.Printf("0. 退出\n")

	var input string

	fmt.Scanln(&input)
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("请输入序号\n")
		SelectFolder()
	}
	number = num
	fmt.Printf("=======================================\n")
	if number == 0 {
		return
	}
	selectReferenceFile()
}

func selectReferenceFile() {
	fmt.Printf("请输入对比文件夹(不输入会默认D盘下同名文件夹):\n")

	fmt.Scanln(&address)
	fmt.Printf("=======================================\n")
	item.NewByNumber(number, address)
	item.RunReference()
	copyToLocal()
	createJSONFile()
}

func copyToLocal() {
	if len(item.SelectFiles) == 0 {
		return
	}
	for _, data := range item.SelectFiles {
		dataName, err := os.Stat(data.Local)
		if err != nil {
			fmt.Println(err)
		}

		item.FileTotal += dataName.Size()
	}
	fmt.Printf("开始复制:")
	service.FileTotal = item.FileTotal
	service.FileSize = item.FileSize

	item.CreateCopy(0)

	fmt.Println("\n完成,大小总共", humanize.Bytes(uint64(item.FileSize)))

	fmt.Printf("=======================================\n")
}

func createJSONFile() {
	fmt.Printf("是否生成JSON\n")
	fmt.Printf("1. 是\n")
	fmt.Printf("2. 否\n")

	var input string

	fmt.Scanln(&input)
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("请输入序号\n")
		SelectFolder()
	}
	item.Number = num

	if item.Number == 2 {
		return
	} else if item.Number == 1 {
		filesData, err := getFilesData(".")
		if err != nil {
			fmt.Println(err)
		}
		b, err := json.Marshal(filesData)
		if err != nil {
			fmt.Println("error:", err)
		}
		//生成json文件
		err = ioutil.WriteFile("test.json", b, 0666)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("=======================================\n")
	}
}

func getFilesData(dirPth string) (*FilesData, error) {
	filesData := new(FilesData)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, fi := range dir {
		if strings.Contains(fi.Name(), "json") {
			break
		} else if fi.IsDir() { // 目录, 递归遍历
			data := FilesData{
				FileName: fi.Name(),
				FileType: "Files",
				Address:  dirPth + "/" + fi.Name(),
			}
			child, err := getFilesData(dirPth + "/" + fi.Name())
			if err != nil {
				fmt.Println(err)
			}
			data.Child = child.Child
			filesData.Child = append(filesData.Child, data)
		} else {
			data := FilesData{
				FileName: fi.Name(),
				FileType: "File",
				Address:  dirPth + "/" + fi.Name(),
			}
			filesData.Child = append(filesData.Child, data)
		}
	}

	return filesData, nil
}
