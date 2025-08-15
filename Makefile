APP := "Sleep guardian"
PKG := ./cmd
OUT := dist
GOFLAGS := -trimpath -ldflags "-s -w"
ICON := internal/tray/sleep_guardian.ico
SYSO := cmd/sleep_guardian.syso

HOSTOS   := $(shell go env GOOS)
HOSTARCH := $(shell go env GOARCH)

TARGETS ?= $(HOSTOS)/$(HOSTARCH)

all: build

clean:
	rm -rf $(OUT)
	rm -f $(SYSO)

build: clean $(TARGETS)
	@echo "Artifacts -> $(OUT)"

$(TARGETS):
	@mkdir -p $(OUT)
	@os="$$(echo $@ | cut -d/ -f1)"; arch="$$(echo $@ | cut -d/ -f2)"; \
	ext=""; if [ "$$os" = "windows" ]; then \
	  ext=".exe"; \
	  go run github.com/akavel/rsrc@latest -ico $(ICON) -o $(SYSO); \
	fi; \
	echo "Building $$os/$$arch"; \
	# darwin требует CGO и тулчейн macOS — пропустим на не-mac \
	if [ "$(HOSTOS)" != "darwin" ] && [ "$$os" = "darwin" ]; then \
	  echo "  -> skip darwin: нужна сборка на macOS"; \
	else \
	  CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(GOFLAGS) -o $(OUT)/$(APP)$$os_$$ext $(PKG); \
	fi; \
	if [ "$$os" = "windows" ]; then rm -f $(SYSO); fi;
	if [ "$$os" = "windows" ]; then \
		copy "configs/config.example.json" "$(OUT)/config.json"; \
	else \
		cp "configs/config.example.json" "$(OUT)/config.json"; \
	fi

fmt: ; go fmt ./...
vet: ; go vet ./...

.PHONY: all clean build fmt vet
