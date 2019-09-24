// pcf - A powerful paste.cf command line client
// Copyright (C) 2019 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"os"
	"path"
	"time"
)

var addr = "paste.cf:21"
var pub = "incoming"
var max = 10 * 1024 * 1024
var stdin_name = "file"

func login() *ftp.ServerConn {
	c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: dial: %v\n", err)
	}
	err = c.Login("anonymous", "anonymous")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: login: %v\n", err)
	}
	return c
}

func store(f *os.File, c *ftp.ServerConn, n string) {
	err := c.Stor(path.Join(pub, n), f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: put: %v\n", err)
	}
}

func exit(c *ftp.ServerConn) {
	err := c.Quit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: quit: %v\n", err)
	}
}

func put(f *os.File, n string) {
	c := login()
	store(f, c, n)
	exit(c)
}

func main() {
	files := os.Args[1:]
	if len(files) == 0 {
		// use stdin data
		put(os.Stdin, stdin_name)
	} else {
		// loop through and use all arguments
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pcf: %v\n", err)
				continue
			}
			put(f, arg)
			f.Close()
		}
	}
}
