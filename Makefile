install: picdown conf scripts/picdown.timer scripts/picdown.service
	install -d /opt/picdown/bin/
	install -d /opt/picdown/conf
	install conf/* /opt/picdown/conf
	install -m 755 picdown /opt/picdown/bin/picdown
	install -d /usr/lib/systemd/system
	install -m 644 scripts/picdown.timer /usr/lib/systemd/system/
	install -m 644 scripts/picdown.service /usr/lib/systemd/system/
	systemctl daemon-reload
	systemctl enable --now picdown.timer picdown.service

build:
	go build -ldflags "-w -s"


.PHONY: install clean