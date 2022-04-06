# pcf - A command line sha1/(S)FTP-based pastebin client.
# Copyright (C) 2022 Dakota Walsh
# GPL3+ See LICENSE in this repo for details.
.POSIX:

include config.mk

all: clean build

build:
	go build
	scdoc < pcf.1.scd | sed "s/VERSION/$(VERSION)/g" > pcf.1

clean:
	rm -f pcf
	rm -f pcf.1

install: build
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f pcf $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/pcf
	mkdir -p $(DESTDIR)$(MANPREFIX)/man1
	cp -f pcf.1 $(DESTDIR)$(MANPREFIX)/man1/pcf.1
	chmod 644 $(DESTDIR)$(MANPREFIX)/man1/pcf.1

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/pcf
	rm -f $(DESTDIR)$(MANPREFIX)/man1/pcf.1

.PHONY: all build clean install uninstall
