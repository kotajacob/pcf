pcf
====

A simple paste.cf command line client. Reads file(s) from STDIN or filename as
argument, uploads to a paste.cf server, and prints the url(s) to STDOUT.

Config files are read to retrieve server information. Optionally server
information can be provided via arguments. See pcf(1) for command line client
information and pcf(5) for config file format information.

License
--------

GPL3+ see LICENSE in this repo for more details.

Build
------

Build dependencies  

 * golang
 * make
 * sed
 * scdoc

`make all`

Install
--------

Optionally configure `config.mk` to specify a different install location.  
Defaults to `/usr/local/`

`sudo make install`

Uninstall
----------

`sudo make uninstall`

Resources
----------

[Send patches](https://git-send-email.io) and questions to
[~kota/pcf@lists.sr.ht](https://lists.sr.ht/~kota/pcf).

Bugs & todo here: [~kota/pcf](https://todo.sr.ht/~kota/pcf)
