VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell git log -1 --format=%cI 2>/dev/null || echo unknown)

VERSION_PKG := github.com/widgora/widgora/internal/version
LDFLAGS     := -s -w -X $(VERSION_PKG).Version=$(VERSION) -X $(VERSION_PKG).Commit=$(COMMIT) -X $(VERSION_PKG).Date=$(DATE)
GOTOOL      := go tool -modfile=tools/go.mod
UI_DIST     := internal/webui/dist
ALLOWED     := $(shell grep -v '^\#' .licenses-allowed | paste -sd, -)

DEV_LISTEN  := 127.0.0.1:8080

.PHONY: build web dev test test-go test-web lint lint-go lint-web vuln fmt licenses licenses-go licenses-web clean

build: web
	CGO_ENABLED=0 go build -trimpath -ldflags '$(LDFLAGS)' -o dist/widgora ./cmd/widgora

web: web/node_modules
	pnpm -C web build
	find $(UI_DIST) -mindepth 1 ! -name .gitkeep -delete
	cp -R web/build/. $(UI_DIST)/

web/node_modules: web/package.json web/pnpm-lock.yaml
	pnpm -C web install --frozen-lockfile
	touch $@

# Go server on $(DEV_LISTEN) plus the Vite dev server, which proxies /api and
# /healthz to it. Ctrl-C stops both.
dev: web/node_modules
	@trap 'kill 0' EXIT INT TERM; \
	go run ./cmd/widgora serve --listen $(DEV_LISTEN) --data-dir .dev-data --log-level debug & \
	pnpm -C web dev

test: test-go test-web

test-go:
	go test -race ./...

test-web: web/node_modules
	pnpm -C web test

lint: lint-go lint-web

lint-go:
	golangci-lint run ./...

lint-web: web/node_modules
	pnpm -C web lint
	pnpm -C web check

vuln:
	$(GOTOOL) govulncheck ./...

fmt: web/node_modules
	golangci-lint fmt ./...
	pnpm -C web format

licenses: licenses-go licenses-web

licenses-go:
	$(GOTOOL) go-licenses check ./... --allowed_licenses=$(ALLOWED)

licenses-web: web/node_modules
	pnpm -C web run check:licenses

clean:
	rm -rf dist web/build
	find $(UI_DIST) -mindepth 1 ! -name .gitkeep -delete
