pcf
====

A simple paste.cf command line client. Reads file(s) from STDIN or filename as
argument, uploads to a paste.cf server, and prints the url(s) to STDOUT.

The PCFSERVER environment variable is used to declare server information. An
example in your shellrc would be `export
PCFSERVER='https://paste.example.com:21/incoming'`. The port and path are optional
and will depend on how the pcf server you're using is configured.

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
