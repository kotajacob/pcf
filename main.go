// pcf - A command line sha1/(S)FTP-based pastebin client.
// Copyright (C) 2022 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
)

func main() {
	// Find and read config.
	configPath, err := xdg.ConfigFile("pcf/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed reading config directory: %v\n", err)
		os.Exit(1)
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed loading config file: %v\n", err)
		os.Exit(1)
	}

	// Parse flags and override config.DefaultMode if needed.
	ftpFlag := flag.Bool("f", false, "use anonymous FTP")
	sftpFlag := flag.Bool("s", false, "use anonymous SFTP")
	authFlag := flag.Bool("a", false, "use authenticated SFTP")

	flag.Parse()

	mode := config.DefaultMode
	switch {
	case *authFlag:
		mode = "sftp-auth"
	case *sftpFlag:
		mode = "sftp-anon"
	case *ftpFlag:
		mode = "ftp"
	}

	// Select the correct upload URL.
	var url string
	switch mode {
	case "sftp-auth":
		url = config.SftpAuthURL
	case "sftp-anon":
		url = config.SftpAnonURL
	default:
		url = config.FtpURL
	}

	// Create requests and upload each file.
	filePaths := flag.Args()
	UploadFiles(
		mode,
		url,
		config.SftpUser,
		config.SftpPass,
		config.Output,
		filePaths,
	)
}

// UploadFiles uploads a list of files by path sequentially in the order given.
// After each successful upload, the resulting uploaded URL is printed to
// STDOUT. Any upload errors will cause an error to be printed to Stderr and
// exit with an error code.
func UploadFiles(mode, url, user, pass, output string, paths []string) {
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed opening %s: %v\n", path, err)
			os.Exit(1)
		}

		var buf bytes.Buffer
		tee := io.TeeReader(f, &buf)

		name := filepath.Base(path)
		if name == "" || name == "." {
			name = RandString(5)
		}

		req, err := NewRequest(mode, url, user, pass, name, tee)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed creating upload request for %s: %v\n", path, err)
			os.Exit(1)
		}

		err = req.Upload()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed while uploading %s: %v\n", path, err)
			os.Exit(1)
		}

		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed closing %s: %v\n", path, err)
			os.Exit(1)
		}

		// Success! Print the output URL.
		fmt.Println(HashName(&buf, path, output))
	}
}

// RandString returns a random latin string of a variable length.
func RandString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// HashName combines a prefix, SHA1 hash from an io.Reader, and the file
// extension from a path.
//
// The result is typically a URL indicating where the uploaded file can be
// located.
func HashName(f io.Reader, path string, prefix string) string {
	var b strings.Builder
	b.WriteString(prefix)

	hash := sha1.New()
	if _, err := io.Copy(hash, f); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"failed calculating hash for %s: %v\n",
			path,
			err,
		)
	}
	b.WriteString(fmt.Sprintf("%x", hash.Sum(nil)))

	b.WriteString(Ext(path))
	return b.String()
}

// Ext returns the file name extension used by path.
// Unlike filepath.Ext we return a blank string for hidden files without an
// extension (.gitignore for example is considered to not have an extension).
func Ext(path string) string {
	path = strings.TrimPrefix(path, ".")
	return filepath.Ext(path)
}
