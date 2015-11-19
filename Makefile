## makefile
.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

MKDIR_P = mkdir -p
RM_F = rm -f
GO = go

DEFAULT_LOG_DIR = /var/lib/sfinder/sregister/log
DEFAULT_CONF_DIR = /etc/sfinder/sregister
DEFAULT_SERVICES_CONF_DIR = /etc/sfinder/sregister/services
DST_BIN_DIR = /usr/bin/
SYSTEMD_FILE_DIR = /usr/lib/systemd/system
SYSTEMD_CONF_DIR = /etc/systemd/system/multi-user.target.wants

EXEC_FILE_NAME = sregister
DEFAULT_CONF_FILE = example/sregister.conf
SYSTEMD_UNIT_FILE = sregister.service

all: test build


build: clean
	@go build $(GOFLAGS)

install:
	${MKDIR_P} ${DEFAULT_LOG_DIR}
	${MKDIR_P} ${DEFAULT_CONF_DIR}
	${MKDIR_P} ${DEFAULT_SERVICES_CONF_DIR}
	${RM_F} ${DEFAULT_CONF_DIR}/${DEFAULT_CONF_FILE}
	${RM_F} ${SYSTEMD_FILE_DIR}/${SYSTEMD_UNIT_FILE}
	${RM_F} ${SYSTEMD_CONF_DIR}/${SYSTEMD_UNIT_FILE}
	@cp ${DEFAULT_CONF_FILE} ${DEFAULT_CONF_DIR}
	@cp ${SYSTEMD_UNIT_FILE} ${SYSTEMD_FILE_DIR}
	@ln -s ${SYSTEMD_FILE_DIR}/${SYSTEMD_UNIT_FILE} ${SYSTEMD_CONF_DIR}/${SYSTEMD_UNIT_FILE}
	@install -d ${DST_BIN_DIR}
	@install -p -D -m 0755 ./${EXEC_FILE_NAME} ${DST_BIN_DIR}  

test:
	@go test $(GOFLAGS) ./configuration

bench: 
	@go test -run=NONE -bench=. $(GOFLAGS) ./...

clean:
	@go clean $(GOFLAGS) -i ./...

## EOF
