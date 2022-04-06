// pcf - A command line sha1/(S)FTP-based pastebin client.
// Copyright (C) 2022 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"fmt"
	"io"
	neturl "net/url"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// A Request to upload a file.
type Request struct {
	// Mode describes the type of request to be made.
	Mode RequestMode

	// URL encodes the upload server, path, username, and password. The path is
	// used to change the directory where the file will be written on the
	// server.
	URL *neturl.URL

	// Name to store the file as on the remote server. This name is overwritten
	// in all modes except RequestModeSFTPAuth. For these other modes you should
	// use a random string in order to not have your file overwritten if another
	// user is simultaneously uploading.
	Name string

	// Content is the data to be uploaded.
	Content io.Reader
}

// RequestMode describes the type of request.
type RequestMode uint8

const (
	// RequestModeFTP is a RequestMode that uploads to an FTP server with an
	// username (defaulting to anonymous) and no password.
	RequestModeFTP RequestMode = iota

	// RequestModeSFTPAnon is a RequestMode that uploads to an SFTP server with
	// an username (defaulting to anonymous) and no password.
	RequestModeSFTPAnon

	// RequestModeSFTPAuth is a RequestMode that uploads to an SFTP server with
	// a username and password.
	RequestModeSFTPAuth
)

// NewRequest returns a new Request given a mode, url, file name, an optional
// username and password and an io.Reader.
//
// Mode is a string representation of a request mode:
// "ftp"       = RequestModeFTP
// "sftp-anon" = RequestModeSFTPAnon
// "sftp-auth" = RequestModeSFTPAuth
// If an invalid mode is provided "ftp" is used.
//
// Usernames and passwords encoded within the URL will be used for any of the
// three modes. However, if user or pass strings are given and the mode is
// "sftp-auth" they will override the encoded credentials. For all other modes
// given credentials are ignored.
func NewRequest(mode string, url string, user string, pass string, name string, content io.Reader) (*Request, error) {
	r := &Request{}
	var err error

	// Set mode and parse url.
	switch mode {
	case "sftp-auth":
		r.Mode = RequestModeSFTPAuth
	case "sftp-anon":
		r.Mode = RequestModeSFTPAnon
	default:
		r.Mode = RequestModeFTP
	}
	r.URL, err = neturl.Parse(url)

	// Replace URL encoded username and password if provided.
	if r.Mode == RequestModeSFTPAuth && user != "" {
		r.URL.User = neturl.UserPassword(user, pass)
	}

	r.Name = name
	r.Content = content
	return r, err
}

// Upload processes a request by uploading it to the FTP or SFTP server and then
// returns nil for a successful upload; or an error.
func (req *Request) Upload() error {
	if req.Mode == RequestModeFTP {
		return req.uploadFTP()
	}
	return req.uploadSFTP()
}

// uploadFTP processes a request using FTP to upload the file.
func (req *Request) uploadFTP() error {
	conn, err := ftp.Dial(req.URL.Host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return fmt.Errorf("connecting to %v: %w", req.URL.Host, err)
	}

	// Check if a password was provided in the URL. If so, use it instead of the
	// default username and password.
	var pass, user string
	if urlPass, ok := req.URL.User.Password(); ok {
		pass = urlPass
		user = req.URL.User.Username()
	} else {
		pass = "anonymous"
		user = "anonymous"
	}

	if err := conn.Login(user, pass); err != nil {
		return fmt.Errorf("logging into anonymous ftp with %s %s: %v", user, pass, err)
	}

	// Store the content, in the path provided in the URL.
	return conn.Stor(filepath.Join(req.URL.Path, req.Name), req.Content)
}

func (req *Request) uploadSFTP() error {
	// Default anonymous username and blank password.
	user := "anonymous"
	pass := ""

	// Check if a password was provided in the URL.
	urlPass, ok := req.URL.User.Password()
	if !ok {
		if req.Mode == RequestModeSFTPAuth {
			return fmt.Errorf("missing sftp password for authenticated sftp mode")
		}
	} else {
		user = req.URL.User.Username()
		pass = urlPass
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: store hostkeys
	}

	// An ssh connection is needed to start sftp.
	conn, err := ssh.Dial("tcp", req.URL.Host, config)
	if err != nil {
		return fmt.Errorf("connecting to %v: %w", req.URL.Host, err)
	}
	defer conn.Close()

	// Start sftp and then upload the file to the given path and name.
	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("failed initializing sftp in ssh session: %v", err)
	}

	path := sftp.Join(req.URL.Path, req.Name)
	f, err := client.Create(path)
	if err != nil {
		return fmt.Errorf("failed creating file at path %s: %v", path, err)
	}

	_, err = f.ReadFrom(req.Content)
	if err != nil {
		return fmt.Errorf("failed writing content to file in sftp connection: %v", err)
	}
	return client.Close()
}
