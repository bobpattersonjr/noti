MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

export GOFLAGS := -mod=vendor
export GOPROXY := off

branch := $(shell git rev-parse --abbrev-ref HEAD)
tag := $(shell git describe --abbrev=0 --tags)
rev := $(shell git rev-parse --short HEAD)

ld_flags_dev := -race -ldflags "-X github.com/bobpattersonjr/noti/internal/command.Version=$(branch)-$(rev)"
ld_flags_rel := -ldflags "-s -w -X github.com/bobpattersonjr/noti/internal/command.Version=$(tag)"

go_src := $(shell find ./service ./internal ./cmd -name "*.go")

go.sum: go.mod
	GOPROXY= go mod tidy

vendor: go.mod go.sum
	GOPROXY= go mod vendor
	touch $@

out/noti: go.mod go.sum vendor $(go_src)
	cd cmd/noti && go build -o ../../$@ $(ld_flags_dev)

out/noti.darwin.rel: go.mod go.sum vendor $(go_src)
	cd cmd/noti && GOOS=darwin GOARCH=amd64 \
		go build -o ../../$@ $(ld_flags_rel)
out/noti.darwin.arm64.rel: go.mod go.sum vendor $(go_src)
	cd cmd/noti && GOOS=darwin GOARCH=arm64 \
		go build -o ../../$@ $(ld_flags_rel)
out/noti.%.rel: go.mod go.sum vendor $(go_src)
	cd cmd/noti && CGO_ENABLED=0 GOOS=$* GOARCH=amd64 \
		go build -o ../../$@ $(ld_flags_rel)

out/noti$(tag).windows-amd64.tar.gz: out/noti.windows.rel
	# Package Windows binary as noti.exe in tarball in a portable way (BSD/GNU tar)
	@tmpdir=$$(mktemp -d); \
	  cp $< $$tmpdir/noti.exe; \
	  tar -czf $@ -C $$tmpdir noti.exe; \
	  rm -rf $$tmpdir
out/noti$(tag).%-amd64.tar.gz: out/noti.%.rel
	# Package binary as noti in tarball in a portable way (BSD/GNU tar)
	@tmpdir=$$(mktemp -d); \
	  cp $< $$tmpdir/noti; \
	  tar -czf $@ -C $$tmpdir noti; \
	  rm -rf $$tmpdir
out/noti$(tag).darwin-arm64.tar.gz: out/noti.darwin.arm64.rel
	# Package darwin arm64 binary in a portable way (BSD/GNU tar)
	@tmpdir=$$(mktemp -d); \
	  cp $< $$tmpdir/noti; \
	  tar -czf $@ -C $$tmpdir noti; \
	  rm -rf $$tmpdir

docs/man/dist/noti.1: docs/man/noti.1.md
	mkdir --parents $(dir $@)
	pandoc -s -t man $< -o $@
docs/man/dist/noti.yaml.5: docs/man/noti.yaml.5.md
	mkdir --parents $(dir $@)
	pandoc -s -t man $< -o $@

.PHONY: build
build: out/noti


.PHONY: lint
lint: goos := $(strip $(shell go env GOOS))
lint: goarch := $(strip $(shell go env GOARCH))
lint: golangci_local := $(if $(filter windows,$(goos)),./tools/golangci-lint-1.64.6-windows-amd64.exe,./tools/golangci-lint-1.64.6-$(goos)-amd64)
lint: golangci_installed := $(shell go env GOBIN)/golangci-lint
lint:
	# Seems like there's some Windows bug with gofmt
	go vet ./...
	@if [ -x "$(golangci_local)" ] && [ "$(goarch)" = "amd64" ]; then \
		"$(golangci_local)" run --no-config --exclude-use-default=false \
		--max-same-issues=0 \
		--timeout 60s \
		--disable errcheck \
		--disable stylecheck \
		--disable bodyclose \
		--$(if $(filter windows,$(goos)),disable,enable) gofmt \
		--$(if $(filter windows,$(goos)),disable,enable) goimports \
		--enable unconvert \
		--enable dupl \
		--enable gocyclo \
		--enable misspell \
		--enable lll \
		--enable unparam \
		--enable nakedret \
		--enable prealloc \
		--enable gocritic \
		--enable gochecknoinits \
		./... ; \
	elif [ -n "$(shell which go 2>/dev/null)" ]; then \
		echo "Installing golangci-lint for $(goos)/$(goarch)..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6; \
		"$(golangci_installed)" run --no-config --exclude-use-default=false \
		--max-same-issues=0 \
		--timeout 60s \
		--disable errcheck \
		--disable stylecheck \
		--disable bodyclose \
		--$(if $(filter windows,$(goos)),disable,enable) gofmt \
		--$(if $(filter windows,$(goos)),disable,enable) goimports \
		--enable unconvert \
		--enable dupl \
		--enable gocyclo \
		--enable misspell \
		--enable lll \
		--enable unparam \
		--enable nakedret \
		--enable prealloc \
		--enable gocritic \
		--enable gochecknoinits \
		./... ; \
	else \
		echo "Skipping golangci-lint: no suitable binary and unable to install"; \
	fi

.PHONY: test
test:
	go test -v -cover -race $(shell go list ./... | grep -v "noti/integration")

.PHONY: test-integration
test-integration: out/noti
	go test -v -cover ./integration/...

.PHONY: release-no-cgo
release-no-cgo: out/noti$(tag).linux-amd64.tar.gz out/noti$(tag).windows-amd64.tar.gz

.PHONY: release-darwin
release-darwin: out/noti$(tag).darwin-amd64.tar.gz out/noti$(tag).darwin-arm64.tar.gz

.PHONY: release-darwin-arm64
release-darwin-arm64: out/noti$(tag).darwin-arm64.tar.gz

.PHONY: release-linux
release-linux: out/noti$(tag).linux-amd64.tar.gz

.PHONY: release-windows
release-windows: out/noti$(tag).windows-amd64.tar.gz

.PHONY: man
man: docs/man/dist/noti.1 docs/man/dist/noti.yaml.5

.PHONY: update-go-mod
update-go-mod:
	GOPROXY= go get -u ./service/... ./internal/... ./cmd/...
