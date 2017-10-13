VERSION := 0.0.1
ROOTDIR := $(shell pwd)
APPDIR := friday
SRCDIR := src/friday
TARGET := bin/friday

GOENV := GOPATH=${ROOTDIR}:${GOPATH} GO15VENDOREXPERIMENT=1

GO := ${GOENV} go
GLIDE := ${GOENV} glide

GLIDEYAML := ${ROOTDIR}/glide.yaml
GLIDELOCK := ${ROOTDIR}/glide.lock

LDFLAGS := -X ${APPDIR}/config.Version=${VERSION}
DEBUGLDFLAGS := ${LDFLAGS} -X ${APPDIR}/config.Mode=debug
RELEASELDFLAGS := -w ${LDFLAGS} -X ${APPDIR}/config.Mode=release

.PHONY: release
release: ${SRCDIR} ${GLIDELOCK}
	${GO} build -i -ldflags="${RELEASELDFLAGS}" -o ${TARGET} friday

.PHONY: build
build: ${SRCDIR} ${GLIDELOCK}
	${GO} build -i -ldflags="${DEBUGLDFLAGS}" -o ${TARGET} friday

${SRCDIR}:
	mkdir -p bin
	mkdir -p src
	ln -s ${ROOTDIR} src/

${GLIDEYAML}:
	${GLIDE} init

${GLIDELOCK}: ${SRCDIR} ${GLIDEYAML}
	${GLIDE} install
	touch ${GLIDELOCK}

.PHONY: init
init: ${SRCDIR} ${GLIDEYAML}

.PHONY: update
update: ${SRCDIR}
	${GLIDE} update

.PHONY: install
install: ${GLIDELOCK}

.PHONY: test
test:
	find "." -name "*_test.go" -not -path "./vendor/*" -not -path "./src/*" -exec dirname {} \; | uniq | xargs go test