pcf(1) [""["pcf-VERSION"]]

# NAME

pcf - Upload files to paste.cf style paste server

# SYNOPSIS

pcf [files]

# DESCRIPTION

pcf reads either a list of files as arguments or data from stdin. Then it
uploads it with ftp and prints the resulting url to stdout.

Set the PCFSERVER environment variable for whatever pcf server you'd like to
use. The format looks like PCFSERVER=https://paste.example.com:21/incoming with
the port and path being optional depending on how the ftp server is configured.

# AUTHORS

Maintained by Dakota Walsh <kota at nilsu.org>.
Up-to-date sources can be found at https://git.sr.ht/~kota/pcf

# COPYRIGHT

Copyright 2019 Dakota Walsh

License GPLv3+ GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
This is free software: you are free to change and redistribute it. There is NO
WARRANTY, to the extent permitted by law.
