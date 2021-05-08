package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type WriteCounter struct {
	Total int64
	Item  int
}

type Option struct {
	Id      int
	Address string
	Local   string
	Status  bool
}

type CopyFile struct {
	OptionData  []Option
	SelectFiles []Option
	Number      int
	Address     string
	FileTotal   int64
	FileSize    int64
}

var lock sync.Mutex

func (data *CopyFile) New() {
	data.OptionData = append(data.OptionData,
		Option{Id: 1, Address: "D:/Animations", Local: "./Animations", Status: true},
		Option{Id: 2, Address: "D:/Flicks", Local: "./Flicks", Status: true},
		Option{Id: 3, Address: "D:/CG", Local: "./CG", Status: true})
}

func (data *CopyFile) NewByNumber(num int, address string) {
	data.OptionData = append(data.OptionData,
		Option{Id: 1, Address: "D:/Animations", Local: "./Animations", Status: true},
		Option{Id: 2, Address: "D:/Flicks", Local: "./Flicks", Status: true},
		Option{Id: 3, Address: "D:/CG", Local: "./CG", Status: true})
	for _, i := range data.OptionData {
		if i.Id != num {
			i.Status = false
		} else {
			if address != "" {
				i.Address = address
			}
		}
	}
}

func CheckDirectory(option []Option) {
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

func (data *CopyFile) RunAutomaticReference() {
	var referenceFiles []Option
	var referenceDirs []Option
	var localFiles []string
	var localDirs []string

	for _, index := range data.OptionData {
		if index.Status {
			file, dirs, _ := getFilesAndDirs(index.Address)
			for _, i := range file {
				rFiles := Option{
					Id:      0,
					Address: i,
					Local:   index.Local,
					Status:  false,
				}

				referenceFiles = append(referenceFiles, rFiles)
			}
			for _, i := range dirs {
				rDirs := Option{
					Id:      0,
					Address: i,
					Local:   index.Local,
					Status:  false,
				}
				referenceDirs = append(referenceDirs, rDirs)
			}

			file, dirs, _ = getFilesAndDirs(index.Local)
			localFiles = append(localFiles, file...)
			localDirs = append(localDirs, dirs...)
		}
	}

	for _, localTable := range localDirs {
		temp, _, _ := getFilesAndDirs(localTable)
		for _, temp1 := range temp {
			localFiles = append(localFiles, temp1)
		}
	}

	for _, referenceTable := range referenceDirs {
		temp, _, _ := getFilesAndDirs(referenceTable.Address)
		for _, temp1 := range temp {
			rFiles := Option{
				Id:      0,
				Address: temp1,
				Local:   referenceTable.Local,
				Status:  false,
			}
			referenceFiles = append(referenceFiles, rFiles)
		}
	}

	for _, table1 := range referenceFiles {
		flag := false
		for _, table2 := range localFiles {
			data1, err := os.Stat(table1.Address)
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
			data.SelectFiles = append(data.SelectFiles, table1)
		}
	}

	for _, selectFile := range data.SelectFiles {
		fmt.Printf("缺少[%s]\n", selectFile.Address)
	}
}

func (data *CopyFile) RunReference() {
	var referenceFiles []Option
	var referenceDirs []Option
	var localFiles []string
	var localDirs []string

	for _, index := range data.OptionData {
		if index.Status {
			file, dirs, _ := getFilesAndDirs(index.Address)
			for _, i := range file {
				rFiles := Option{
					Id:      0,
					Address: i,
					Local:   index.Local,
					Status:  false,
				}

				referenceFiles = append(referenceFiles, rFiles)
			}
			for _, i := range dirs {
				rDirs := Option{
					Id:      0,
					Address: i,
					Local:   index.Local,
					Status:  false,
				}
				referenceDirs = append(referenceDirs, rDirs)
			}

			file, dirs, _ = getFilesAndDirs(index.Local)
			localFiles = append(localFiles, file...)
			localDirs = append(localDirs, dirs...)
		}
	}

	for _, localTable := range localDirs {
		temp, _, _ := getFilesAndDirs(localTable)
		for _, temp1 := range temp {
			localFiles = append(localFiles, temp1)
		}
	}

	for _, referenceTable := range referenceDirs {
		temp, _, _ := getFilesAndDirs(referenceTable.Address)
		for _, temp1 := range temp {
			rFiles := Option{
				Id:      0,
				Address: temp1,
				Local:   referenceTable.Local,
				Status:  false,
			}
			referenceFiles = append(referenceFiles, rFiles)
		}
	}

	for _, table1 := range referenceFiles {
		flag := false
		for _, table2 := range localFiles {
			data1, err := os.Stat(table1.Address)
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
			data.SelectFiles = append(data.SelectFiles, table1)
		}
	}

	for _, selectFile := range data.SelectFiles {
		fmt.Printf("缺少[%s]\n", selectFile.Address)
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
