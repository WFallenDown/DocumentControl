package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

var number uint64

var dir string

var n int

func main() {
	var over string
	fmt.Printf("test")
	fmt.Println("\rtest2")

	dir = filepath.Dir(os.Args[0])
	var c []chan int

	for i := 0; i < 10; i++ {
		c = append(c, make(chan int))
		go runTest(i, c[i])
	}

	for _, index := range c {
		d := <-index
		fmt.Println(d)
	}
	fmt.Println(dir)
	fmt.Scanln(&over)
}

func runTest(i int, c chan int) {
	//time.Sleep(time.Duration(i) * time.Second)
	n = i
	fmt.Printf("--------------\n")
	fmt.Println("n:", n, "i:", i)
	fmt.Printf("..............\n")
	c <- i
	close(c)
}

func run(n int, c chan int) {
	srcFile, err := os.Open("/Users/wangbin/OneDrive/Windows.iso")
	if err != nil {
		fmt.Printf("open file err=%v\n", err)
	}
	defer srcFile.Close()
	//通过src file ,获取到 Reader
	reader := bufio.NewReader(srcFile)
	//打开dstFileName
	dstFile, err := os.OpenFile(dir+"/Windows"+strconv.Itoa(n)+".iso", os.O_RDWR|os.O_CREATE, 0666)
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
	Reader, err := io.Copy(writer, io.TeeReader(reader, counter))

	// If error is not nil then panics
	if err != nil {
		panic(err)
	}

	c <- 1

	// Prints output
	fmt.Printf("n:%v\n", Reader)
}

type WriteCounter struct {
	Total int64
	Item  int
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += int64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc *WriteCounter) PrintProgress() {
	num := float64(wc.Total) / float64(number)
	f := int(math.Floor((num * 100) + 0.5))

	fmt.Printf("\r %d %%", f)
}
