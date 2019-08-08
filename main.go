// pcf - A powerful paste.cf command line client
// Copyright (C) 2019 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"bufio"
	"fmt"
	"os"
)

// var addr = "paste.cf"
// var pub = "incoming"
// var max = 10 * 1024 * 1024
func upload(f *os.File) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		fmt.Println(input.Text())
	}
}

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		upload(os.Stdin)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pcf: %v\n", err)
				continue
			}
			upload(f)
			f.Close()
		}
	}
	fmt.Println("Done")
}
