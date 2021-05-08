package service

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

var FileTotal int64
var FileSize int64

func CreateCopy(data *Option, c chan bool) {

	_, err := copyFile(*data)
	if err != nil {
		fmt.Println(err)
		c <- false
	}
	c <- true
	close(c)
}

func copyFile(path Option) (written int64, err error) {
	dstFileName := path.Local

	files := strings.Split(path.Address, "/")

	verifyFolder := dstFileName + "/" + files[len(files)-2] + "/" + files[len(files)-1]

	//pathExists(verifyFolder)

	srcFile, err := os.Open(path.Address)
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

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total = int64(n)

	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	lock.Lock()
	FileSize += wc.Total
	num := float64(FileSize) / float64(FileTotal)
	f := int(math.Floor((num * 100) + 0.5))

	fmt.Printf("\r %d %%", f)
	lock.Unlock()
}
