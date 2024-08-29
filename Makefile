# get the repo root and output path
ROOT_PACKAGE=github.com/biomed168/sdk
OUT_DIR=$(REPO_ROOT)/_output

# define the default goal
SHELL := /bin/bash
DIRS=$(shell ls)
GO=go

# 设置默认 shell 为 /bin/bash，并通过 ls 命令获取当前目录下的所有文件夹和文件
.DEFAULT_GOAL := help

#获取当前 Makefile 文件的目录
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

# 设置代码库的根目录，若 ROOT_DIR 未定义，则将其设置为当前 Makefile 文件所在目录的绝对路径
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/. && pwd -P))
endif

# 定义输出目录
# 设置构建输出目录 _output，并确保该目录存在
ifeq ($(origin OUTPUT_DIR),undefined)
OUTPUT_DIR := $(ROOT_DIR)/_output
$(shell mkdir -p $(OUTPUT_DIR))
endif

# 设置二进制文件输出目录 _output/bin，并确保该目录存在
ifeq ($(origin BIN_DIR),undefined)
BIN_DIR := $(OUTPUT_DIR)/bin
$(shell mkdir -p $(BIN_DIR))
endif

# 设置工具输出目录 _output/tools，并确保该目录存在
ifeq ($(origin TOOLS_DIR),undefined)
TOOLS_DIR := $(OUTPUT_DIR)/tools
$(shell mkdir -p $(TOOLS_DIR))
endif

# 设置临时文件目录 _output/tmp，并确保该目录存在
ifeq ($(origin TMP_DIR),undefined)
TMP_DIR := $(OUTPUT_DIR)/tmp
$(shell mkdir -p $(TMP_DIR))
endif

# 版本信息
# 从 Git 标签和提交记录中获取版本号，默认格式为 v2.3.3.631.g00abdc9b.dirty
ifeq ($(origin VERSION), undefined)
VERSION := $(shell git describe --tags --always --match="v*" --dirty | sed 's/-/./g')
endif

# 获取当前 Git 的状态（是否有未提交的修改），并获取当前提交的哈希值
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

# 设置构建目标文件和输出路径
BUILDFILE = "./main.go"
BUILDAPP = "$(OUTPUT_DIR)/"

# 代码目录定义
# 定义需要扫描的代码目录，并设置 find 命令的基本路径
CODE_DIRS := $(ROOT_DIR)/
FINDS := find $(CODE_DIRS)

# 平台相关配置
# 定义支持的构建平台列表
PLATFORMS ?= darwin_amd64 darwin_arm64 windows_amd64 linux_amd64 linux_arm64

# 根据当前环境变量或者指定平台来设置构建的操作系统和架构
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
		GOARCH := $(shell go env GOARCH)
	endif
	PLATFORM := $(GOOS)_$(GOARCH)
	IMAGE_PLAT := linux_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif

# 定义 find 和 xargs 命令，过滤掉不需要扫描的路径
FIND := find . ! -path './image/*' ! -path './vendor/*' ! -path './bin/*'
XARGS := xargs -r

# 定义支持的 Go 版本
GO_SUPPORTED_VERSIONS ?= 1.18|1.19|1.20|1.21

# 在构建过程中通过 -ldflags 传递版本信息
GO_LDFLAGS += -X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# 设置 Go 构建标志，支持调试时禁用优化
ifeq ($(DLV),)
	GO_BUILD_FLAGS += -gcflags "all=-N -l"
	LDFLAGS = ""
endif
GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

# 针对 Windows 平台，生成 .exe 扩展名的可执行文件
ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif


# GOPATH 和 GOBIN 设置
# 获取 Go 的 GOPATH 并设置 GOBIN，用于存放安装的工具

GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

# 获取 cmd 目录下所有子目录中的命令，并从路径中提取二进制文件名

ifeq (${COMMANDS},)
  $(error Could not determine COMMANDS, set ROOT_DIR or run in source dir)
endif
ifeq (${BINS},)
  $(error Could not determine BINS, set ROOT_DIR or run in source dir)
endif

# 如果没有找到命令或二进制文件，输出错误信息
EXCLUDE_TESTS=github.com/biomed168/sdk/test


# 构建
.PHONY: all
all: build

# 定义 build 目标，根据当前平台进行构建，并生成对应的二进制文件
.PHONY: build
build:
	@echo "===========> Building for $(OS)/$(ARCH)"
	@CGO_ENABLED=1 GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BIN_DIR)/accompany-sdk-core-$(OS)-$(ARCH) $(TARGET)

.PHONY: build-multiple
build-multiple:
	@for os in $(OSES); do \
		for arch in $(ARCHS); do \
			$(MAKE) build OS=$$os ARCH=$$arch; \
		done \
	done

# 定义 build-multiple 目标，用于构建所有支持的平台

.PHONY: build-wasm
build-wasm:
	GOOS=js GOARCH=wasm go build -trimpath -ldflags "-s -w" -o ${BIN_DIR}/accompany.wasm wasm/cmd/main.go

## install: Install the binary to the BIN_DIR
.PHONY: install
install: build
	mv ${BINARY_NAME} ${BIN_DIR}

## ios: Build the iOS framework
.PHONY: ios
ios:
	go get golang.org/x/mobile
	rm -rf build/ accompany_sdk/t_friend_sdk.go accompany_sdk/t_group_sdk.go  accompany_sdk/ws_wrapper/
	GOARCH=arm64 gomobile bind -v -trimpath -ldflags "-s -w" -o build/accompanyCore.xcframework -target=ios ./accompany_sdk/ ./accompany_sdk_callback/

.PHONY: android
android:
	go get golang.org/x/mobile/bind
	GOARCH=amd64 gomobile bind -v -trimpath -ldflags="-s -w" -o ./accompany_sdk.aar -target=android ./accompany_sdk/ ./accompany_sdk_callback/

.PHONY: test
test:
	$(GO) test $(shell $(FIND) -name '*_test.go' | sed 's|/[^/]*$$||' | sort | uniq | sed 's|^\./|./|') $(TESTFLAGS)

.PHONY: clean
clean:
	rm -vrf $(OUTPUT_DIR)/*
	go clean --modcache

.PHONY: tools.verify.%
tools.verify.%:
	@$(TOOLS_DIR)/$* --version > /dev/null 2>&1 || { echo "installing $*"; $(MAKE) install.$*; }

.PHONY: help
help: Makefile
	@sed -n 's/^##//p' $<


















