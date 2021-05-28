// pcf - A powerful paste.cf command line client
// Copyright (C) 2019 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

var addr = os.Getenv("PCFSERVER")

// create the ftp connection
func login(u *url.URL) *ftp.ServerConn {
	c, err := ftp.Dial(u.Host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: dial: %v\n", err)
	}
	err = c.Login("anonymous", "anonymous")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: login: %v\n", err)
	}
	return c
}

// store the file in the connection
func store(u *url.URL, f *os.File, c *ftp.ServerConn, n string) {
	err := c.Stor(path.Join(u.Path, n), f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: put: %v\n", err)
	}
}

// close the ftp connection
func exit(c *ftp.ServerConn) {
	err := c.Quit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: quit: %v\n", err)
	}
}

// upload the file to the ftp server
func put(f *os.File, n string, u *url.URL) {
	if _, err := f.Seek(0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "pcf: failed to read: %v\n", err)
		os.Exit(1)
	}
	c := login(u)
	store(u, f, c, n)
	exit(c)
}

// calculate and print the hash
func hash(f *os.File) string {
	if _, err := f.Seek(0, 0); err != nil {
		fmt.Fprintf(os.Stderr, "pcf: failed to read: %v\n", err)
		os.Exit(1)
	}
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Fprintf(os.Stderr, "pcf: failed to hash: %v\n", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func main() {
	if addr == "" {
		fmt.Println("pcf: you must set the PCFSERVER environment variable!")
		os.Exit(1)
	}
	// parse the url
	u, err := url.Parse(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: url configuration wrong: %v\n", err)
	}
	args := os.Args[1:]
	if len(args) == 0 {
		// use stdin data
		put(os.Stdin, "file", u)
		hash(os.Stdin)
	} else {
		// loop through and use all arguments
		for _, arg := range args {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pcf: open: %v\n", err)
				continue
			}
			defer f.Close()

			// upload the file
			put(f, filepath.Base(arg), u)
			// calculate the hash
			h := hash(f)
			// print the url
			u.Host = u.Hostname()
			u.Path = h + filepath.Ext(arg)
			fmt.Println(u)
		}
	}
}
