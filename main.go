// pcf - A powerful paste.cf command line client
// Copyright (C) 2019 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

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
	addr := os.Getenv("PCFSERVER")
	if addr == "" {
		fmt.Println("pcf: you must set the PCFSERVER environment variable!")
		os.Exit(1)
	}

	// parse the ftpURL
	ftpURL, err := url.Parse(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: url configuration wrong: %v\n", err)
	}
	args := os.Args[1:]

	if len(args) == 0 {
		// use stdin data
		c := login(ftpURL)
		store(ftpURL, os.Stdin, c, "file")
		exit(c)
		h := hash(os.Stdin)
		webURL := *ftpURL
		webURL.Host = ftpURL.Hostname()
		webURL.Path = h
		fmt.Println(&webURL)
	} else {
		// loop through and use all arguments
		connection := login(ftpURL)
		for _, arg := range args {
			file, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "pcf: open: %v\n", err)
				continue
			}
			defer file.Close()

			if _, err := file.Seek(0, 0); err != nil {
				fmt.Fprintf(os.Stderr, "pcf: failed to read: %v\n", err)
				os.Exit(1)
			}
			store(ftpURL, file, connection, filepath.Base(arg))

			// calculate the hash
			h := hash(file)

			// print the URL
			webURL := *ftpURL
			webURL.Host = ftpURL.Hostname()
			webURL.Path = h + filepath.Ext(arg)
			fmt.Println(&webURL)
		}
		exit(connection)
	}
}

// create the ftp connection
func login(u *url.URL) *ftp.ServerConn {
	c, err := ftp.Dial(u.Host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "pcf: %v\n", err)
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
