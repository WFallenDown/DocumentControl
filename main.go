package main

import (
	"DocumentControl/automatic"
	"DocumentControl/manual"
	"fmt"
)

func main() {
	var over string
	fmt.Printf("=======================================\n")
	fmt.Printf("Welcome to file copy\n")
	work := automatic.Run()
	if work == true {
		fmt.Printf("结束")
		fmt.Scanln(&over)
		return
	}
	manual.SelectFolder()
	fmt.Printf("结束")
	fmt.Scanln(&over)
}
