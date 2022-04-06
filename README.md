# pcf
A command line sha1/(S)FTP-based pastebin client. Reads from STDIN or a list of
filenames as arguments, uploads files to a server, and prints the urls to
STDOUT.

Server information is configured using a simple `toml` format in
`$XDG_CONFIG_HOME/pcf/config.toml` which is typically
`$HOME/.config/pcf/config.toml`. An example configuration is shown below with
the details for a free public pcf server
[paste.nilsu.org](https://paste.nilsu.org).

## Config
```toml
ftp_url = "ftp://paste.nilsu.org:21/incoming"
sftp_anon_url = "sftp://paste.nilsu.org:22/incoming"
sftp_auth_url = "sftp://paste.nilsu.org:22/var/www/html/paste"
sftp_user = "kota"
sftp_pass = "cowscows"

# Default mode. Options are "ftp", "sftp-anon", and "sftp-auth".
default_mode = "sftp-anon"

# Prefix for the resulting file sha1 name.
output = "https://paste.nilsu.org/"
```

## License
GPL3+ see LICENSE in this repo for more details.

## Build
Build dependencies:
 * golang
 * make
 * sed
 * scdoc

`make all`

## Install
Optionally configure `config.mk` to specify a different install location.\
Defaults to `/usr/local/`

`sudo make install`

## Uninstall
`sudo make uninstall`

## Resources
[Send patches](https://git-send-email.io) and questions to
[~kota/pcf@lists.sr.ht](https://lists.sr.ht/~kota/pcf).

Bugs & todo here: [~kota/pcf](https://todo.sr.ht/~kota/pcf)
