package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var number uint64

func main() {

	srcFile, err := os.Open("/Users/wangbin/OneDrive/Windows.iso")
	if err != nil {
		fmt.Printf("open file err=%v\n", err)
	}
	defer srcFile.Close()
	//通过src file ,获取到 Reader
	reader := bufio.NewReader(srcFile)

	//打开dstFileName
	dstFile, err := os.OpenFile("./Windows.iso", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("open file err=%v\n", err)
		return
	}

	//通过dstFile, 获取到 Writer
	writer := bufio.NewWriter(dstFile)
	defer dstFile.Close()

	counter := &WriteCounter{}
	data1, err := os.Stat("/Users/wangbin/OneDrive/Windows.iso")
	number = uint64(data1.Size())
	Reader, err:= io.Copy(writer, io.TeeReader(reader, counter))

	// If error is not nil then panics
	if err != nil {
		panic(err)
	}

	// Prints output
	fmt.Printf("n:%v\n", Reader)
}

type WriteCounter struct {
	Total int64
	Item  float64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	fmt.Print(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	num := float64(wc.Total) / float64(number)
	f := num * 100

	if wc.Item == 0.00 {
		fmt.Printf("进度: [")
		wc.Item += 1
	}
	if wc.Item < 100.00 && f >= wc.Item {
		fmt.Printf("=")
		wc.Item += 1
	} else if f >= 100 {
		fmt.Printf("]\n")
	}

}
