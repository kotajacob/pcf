// pcf - A command line sha1/(S)FTP-based pastebin client.
// Copyright (C) 2022 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config a user's options for pcf.
type Config struct {
	// URL for anonymous FTP uploads. Anonymous FTP is an FTP server with an
	// account (defaulting to anonymous) which has no password.
	FtpURL string `toml:"ftp_url"`

	// URL for anonymous SFTP uploads. Anonymous SFTP is an SFTP server with an
	// account (defaulting to anonymous) which has no password.
	SftpAnonURL string `toml:"sftp_anon_url"`

	// URL for authenticated SFTP uploads.
	SftpAuthURL string `toml:"sftp_auth_url"`

	// Username for authenticated SFTP uploads.
	SftpUser string `toml:"sftp_user"`

	// Password for authenticated SFTP uploads.
	SftpPass string `toml:"sftp_pass"`

	// Options are "ftp", "sftp-anon", and "sftp-auth".
	// When using "sftp-auth": SftpUser and SftpPass are used for
	// authentication.
	//
	// If unset "ftp" will be used as the default mode.
	DefaultMode string `toml:"default_mode"`

	// Output is a prefix, typically a url, to be printed before the new
	// filename.
	Output string `toml:"output"`
}

// LoadConfig reads a config from a filename and returns a Config.
func LoadConfig(filename string) (*Config, error) {
	c := &Config{}
	f, err := os.Open(filename)
	if err != nil {
		return c, fmt.Errorf("failed opening config file: %v", err)
	}

	t := toml.NewDecoder(f)
	_, err = t.Decode(c) // TOML metadata is ignored.
	if err != nil {
		return c, fmt.Errorf("failed parsing config file: %v", err)
	}

	if err := f.Close(); err != nil {
		return c, fmt.Errorf("failed closing config file: %v", err)
	}

	return c, err
}
