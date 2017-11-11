VERSION := 0.0.1
ROOTDIR := $(shell pwd)
APPNAME := friday
SRCDIR := src/friday
TARGET := bin/friday

GOENV := GOPATH=${ROOTDIR} GO15VENDOREXPERIMENT=1

GO := ${GOENV} go
GLIDE := ${GOENV} glide

GLIDEYAML := ${ROOTDIR}/glide.yaml
GLIDELOCK := ${ROOTDIR}/glide.lock

LDFLAGS := -X ${APPNAME}/config.Version=${VERSION}
DEBUGLDFLAGS := ${LDFLAGS} -X ${APPNAME}/config.Mode=debug
RELEASELDFLAGS := -w ${LDFLAGS} -X ${APPNAME}/config.Mode=release

DEBUGBUILDTAGS := debug dbsqlite
RELEASEBUILDTAGS := release dball

.PHONY: release
release: ${SRCDIR}
	${GO} build -i -tags "${RELEASEBUILDTAGS}" -ldflags="${RELEASELDFLAGS}" -o ${TARGET} friday

.PHONY: build
build: ${SRCDIR}
	${GO} build -i -tags "${DEBUGBUILDTAGS}" -ldflags="${DEBUGLDFLAGS}" -o ${TARGET} friday

${SRCDIR}:
	mkdir -p bin
	mkdir -p src
	ln -s `dirname "${ROOTDIR}"`/${APPNAME} src/

${GLIDEYAML}:
	${GLIDE} init

${GLIDELOCK}: ${SRCDIR} ${GLIDEYAML}
	${GLIDE} install
	touch ${GLIDELOCK}

.PHONY: init
init: ${SRCDIR} ${GLIDEYAML}

.PHONY: dev-init
dev-init: init
	echo 'make test || exit $?' > .git/hooks/pre-push
	chmod +x .git/hooks/pre-push


.PHONY: update
update: ${SRCDIR}
	${GLIDE} update
	find vendor -name 'testdata' -type d -exec rm -rf {} \; || true
	find vendor -name '*_test.go' -delete || true
	find vendor -type f \( ! -name '*.go' ! -name 'LICENSE' ! -name '*.s' ! -name '*.h' ! -name '*.c' ! -name '*.cpp' \) -delete || true

.PHONY: install
install: ${GLIDELOCK}

.PHONY: test
test: init
	$(eval packages ?= $(patsubst ./%,${APPNAME}/%,$(shell find "." -name "*_test.go" -not -path "./vendor/*" -not -path "./src/*" -not -path "./.*" -exec dirname {} \; | uniq)))
	${GOENV} go test ${packages}

.PHONY: lint
lint:
	${GOENV} find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./src/*" -exec golint {} \;

.PHONY: migration
migration:
	$(eval migration ?= `date +%y%m%d%H%M%S`)
	@echo "package migration \n\
	\n\
	import \"friday/storage\"\n\
	\n\
	// Migrate${migration} :\n\
	func (c *Command) Migrate${migration}(migration *Migration, conn *storage.DatabaseConnection) error { \n\
	\n\
	}\n\
	\n\
	// Rollback${migration} :\n\
	func (c *Command) Rollback${migration}(migration *Migration, conn *storage.DatabaseConnection) error { \n\
	    return nil\n\
	}"> storage/migration/migration${migration}.go