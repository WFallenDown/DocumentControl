package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"
)

type WriteCounter struct {
	Total int64
	Item  int
}

type FilesData struct {
	FileName string
	FileType string
	Tag      []string
	Describe string
	Address  string
	Child    []FilesData
}

type Option struct {
	Id      int
	Address string
	Status  bool
}

var number int
var address string
var selectFiles []string
var over string
var option []Option
var fileTotal int64
var fileSize int64
var lock sync.Mutex

func main() {
	option = append(option, Option{Id: 1, Address: "D:/Animations", Status: true},
		Option{Id: 2, Address: "D:/Flicks", Status: true},
		Option{Id: 3, Address: "D:/CG", Status: true})
	work := automaticOperation()
	if work == true {
		fmt.Printf("结束")
		fmt.Scanln(&over)
		return
	}
	selectFolder()
	fmt.Printf("结束")
	fmt.Scanln(&over)
}

func automaticOperation() bool {
	fmt.Printf("=======================================\n")
	fmt.Printf("自动拷贝?\n")
	fmt.Printf("1. 是\n")
	fmt.Printf("2. 否\n")

	var input string
	fmt.Scanln(&input)
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("请输入序号\n")
		automaticOperation()
	}
	if num == 2 {
		return false
	}

	automaticRun()
	return true
}

func checkDirectory() {
	fmt.Printf("=======================================\n")
	fmt.Printf("检测文件夹是否存在\n")
	for _, index := range option {
		pathExists(&index)
		if index.Status == false {
			fmt.Printf("%v不存在，将跳过此文件夹\n", index.Address)
		}
	}
	fmt.Printf("检测完成")
}

func automaticRun() {

	checkDirectory()
	selectFiles = []string{}
	for _, index := range option {
		if index.Status {
			number = index.Id
			runReference()
		}
	}
	automaticCopy()
	//createJSONFile()
}

func selectFolder() {
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
		selectFolder()
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
	runReference()
	copyToLocal()
	createJSONFile()
}

func runReference() {
	var referenceFiles []string
	var referenceDirs []string
	var localFiles []string
	var localDirs []string
	switch number {
	case 1:
		localFiles, localDirs, _ = getFilesAndDirs("./Animations")
		if address == "" {
			referenceFiles, referenceDirs, _ = getFilesAndDirs("D:/Animations")
		} else {
			referenceFiles, referenceDirs, _ = getFilesAndDirs(address)
		}
		break
	case 2:
		localFiles, localDirs, _ = getFilesAndDirs("./Flicks")
		if address == "" {
			referenceFiles, referenceDirs, _ = getFilesAndDirs("D:/Flicks")
		} else {
			referenceFiles, referenceDirs, _ = getFilesAndDirs(address)
		}
		break
	case 3:
		localFiles, localDirs, _ = getFilesAndDirs("./CG")
		if address == "" {
			referenceFiles, referenceDirs, _ = getFilesAndDirs("D:/CG")
		} else {
			referenceFiles, referenceDirs, _ = getFilesAndDirs(address)
		}
		break
	}

	for _, localTable := range localDirs {
		temp, _, _ := getFilesAndDirs(localTable)
		for _, temp1 := range temp {
			localFiles = append(localFiles, temp1)
		}
	}

	for _, referenceTable := range referenceDirs {
		temp, _, _ := getFilesAndDirs(referenceTable)
		for _, temp1 := range temp {
			referenceFiles = append(referenceFiles, temp1)
		}
	}

	for _, table1 := range referenceFiles {
		flag := false
		for _, table2 := range localFiles {
			data1, err := os.Stat(table1)
			if err != nil {
				fmt.Println(err)
			}
			data2, err := os.Stat(table2)
			if err != nil {
				fmt.Println(err)
			}
			if data1.Name() == data2.Name() && data1.Size() == data2.Size() {
				flag = true
			}
		}
		if !flag {
			selectFiles = append(selectFiles, table1)
		}
	}

	for _, selectFile := range selectFiles {
		fmt.Printf("缺少[%s]\n", selectFile)
	}
}

func automaticCopy() {
	if len(selectFiles) == 0 {
		return
	}
	for _, data := range selectFiles {
		dataName, err := os.Stat(data)
		if err != nil {
			fmt.Println(err)
		}

		fileTotal += dataName.Size()
	}
	fmt.Printf("开始复制:")

	var c []chan bool

	for index, data := range selectFiles {
		c = append(c, make(chan bool))
		go createCopy(data, c[index])
	}
	for _, data := range c {
		if <-data == false {
			return
		}
	}
	fmt.Println("\n完成,大小总共", humanize.Bytes(uint64(fileSize)))

	fmt.Printf("=======================================\n")
}

func copyToLocal() {
	for _, data := range selectFiles {
		dataName, err := os.Stat(data)
		if err != nil {
			fmt.Println(err)
		}

		fileTotal += dataName.Size()
	}
	fmt.Printf("开始复制:")

	var c []chan bool

	for index, data := range selectFiles {
		c = append(c, make(chan bool))
		go createCopy(data, c[index])
	}
	for _, data := range c {
		if <-data == false {
			return
		}
	}
	fmt.Println("\n完成,大小总共", humanize.Bytes(uint64(fileSize)))

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
		selectFolder()
	}
	number = num

	if number == 2 {
		return
	} else if number == 1 {
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

//获取指定目录下的所有文件和目录
func getFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+"/"+fi.Name())
			getFilesAndDirs(dirPth + "/" + fi.Name())
		} else {
			con := []string{"AVI", "mov", "rmvb", "rm", "FLV", "mp4", "3GP"}
			flag := false
			for _, str := range con {
				if strings.Contains(fi.Name(), str) {
					flag = true
					break
				}
			}
			if flag {
				files = append(files, dirPth+"/"+fi.Name())
			}
		}
	}

	return files, dirs, nil
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

func createCopy(data string, c chan bool) {

	_, err := copyFile(data)
	if err != nil {
		fmt.Println(err)
		c <- false
	}
	c <- true
	close(c)
}

func copyFile(srcFileName string) (written int64, err error) {
	var dstFileName string
	if number == 1 {
		dstFileName = "."
	} else if number == 2 {
		dstFileName = "./Flicks"
	}

	files := strings.Split(srcFileName, "/")

	verifyFolder := dstFileName + "/" + files[len(files)-2] + "/" + files[len(files)-1]

	//pathExists(verifyFolder)

	srcFile, err := os.Open(srcFileName)
	if err != nil {
		fmt.Printf("open file err=%v\n", err)
	}
	defer srcFile.Close()
	//通过src file ,获取到 Reader
	reader := bufio.NewReader(srcFile)

	//打开dstFileName
	dstFile, err := os.OpenFile(verifyFolder, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err=%v\n", err)
		return
	}

	//通过dstFile, 获取到 Writer
	writer := bufio.NewWriter(dstFile)
	defer dstFile.Close()

	counter := &WriteCounter{}
	return io.Copy(writer, io.TeeReader(reader, counter))
}

//pathExists 判断文件夹是否存在
func pathExists(data *Option) {
	_, err := os.Stat(data.Address)
	if err != nil {
		data.Status = false
	}
	/*if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			return true, nil
		}
	}*/
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total = int64(n)

	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	lock.Lock()
	fileSize += wc.Total
	num := float64(fileSize) / float64(fileTotal)
	f := int(math.Floor((num * 100) + 0.5))

	fmt.Printf("\r %d %%", f)
	lock.Unlock()
}
