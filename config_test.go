// pcf - A command line sha1/(S)FTP-based pastebin client.
// Copyright (C) 2022 Dakota Walsh
// GPL3+ See LICENSE in this repo for details.
package main

import (
	"os"
	"path/filepath"
	"testing"
)

var (
	testConfigData = []byte(`ftp_url = "ftp://paste.nilsu.org:21/incoming"
sftp_anon_url = "sftp://paste.nilsu.org:22/incoming"
sftp_auth_url = "sftp://paste.nilsu.org:22/var/www/html/paste"
sftp_user = "kota"
sftp_pass = "cowscows"
# Default mode. Options are "ftp" "sftp-anon" and "sftp-auth".
default_mode = "sftp-anon"
output = "https://paste.nilsu.org/"`)
)

func TestLoadConfig(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	os.WriteFile(
		configPath,
		testConfigData,
		0666,
	)

	c, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("failed reading config at %s: %v\n", configPath, err)
	}

	if c.FtpURL != "ftp://paste.nilsu.org:21/incoming" {
		t.Fatal("incorrect ftp_url from testConfigFile")
	}
	if c.SftpAnonURL != "sftp://paste.nilsu.org:22/incoming" {
		t.Fatal("incorrect sftp_anon_url from testConfigFile")
	}
	if c.SftpAuthURL != "sftp://paste.nilsu.org:22/var/www/html/paste" {
		t.Fatal("incorrect sftp_auth_url from testConfigFile")
	}
	if c.SftpUser != "kota" {
		t.Fatal("incorrect sftp_user from testConfigFile")
	}
	if c.SftpPass != "cowscows" {
		t.Fatal("incorrect sftp_pass from testConfigFile")
	}
	if c.DefaultMode != "sftp-anon" {
		t.Fatal("incorrect default_mode from testConfigFile")
	}
	if c.Output != "https://paste.nilsu.org/" {
		t.Fatal("incorrect output from testConfigFile")
	}
}
