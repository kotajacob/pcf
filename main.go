// pcf - A powerful paste.cf command line client
// Copyright (C) 2019 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"os"
	"path"
	"time"
)

var addr = "paste.cf:21"
var pub = "incoming"
var max = 10 * 1024 * 1024

// upload files to a connection
func upload() {
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: dial: %v\n", err)
	}
	err = c.Login("anonymous", "anonymous")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: login: %v\n", err)
	}
	data := bytes.NewBufferString("Hello World\n")
	err = c.Stor(path.Join(pub, "test-file.txt"), data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: put: %v\n", err)
	}
	err = c.Quit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: quit: %v\n", err)
	}
}

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		upload()
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pcf: %v\n", err)
				continue
			}
			upload()
			f.Close()
		}
	}
	fmt.Println("Done")
}
