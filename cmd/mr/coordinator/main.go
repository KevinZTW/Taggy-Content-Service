package main

import (
	"content-parser/pkg/mr"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Parser Master Start")
	files := []string{"../../src/file1.txt", "../../src/file2.txt"}
	m := mr.MakeCoordinator(files, 2)
	for m.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
