VERSION := 0.0.1
ROOTDIR := $(shell pwd)
APPDIR := friday
SRCDIR := src/friday
TARGET := bin/friday

GOENV := GOPATH=${ROOTDIR} GO15VENDOREXPERIMENT=1

GO := ${GOENV} go

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

.PHONY: init
init: ${SRCDIR}

.PHONY: dev-init
dev-init: init
	echo 'make test || exit $?' > .git/hooks/pre-push
	chmod +x .git/hooks/pre-push

.PHONY: update
update: ${SRCDIR}
	cd ${SRCDIR} && ${GOENV} godep save ${APPDIR}

.PHONY: test
test:
	find "." -name "*_test.go" -not -path "./vendor/*" -not -path "./src/*" -exec dirname {} \; | uniq | xargs env ${GOENV} go test

.PHONY: lint
lint:
	${GOENV} find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./src/*" -exec golint {} \;