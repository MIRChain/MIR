# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: mir android ios mir-cross evm all test clean
.PHONY: mir-linux mir-linux-386 mir-linux-amd64 mir-linux-mips64 mir-linux-mips64le
.PHONY: mir-linux-arm mir-linux-arm-5 mir-linux-arm-6 mir-linux-arm-7 mir-linux-arm64
.PHONY: mir-darwin mir-darwin-386 mir-darwin-amd64
.PHONY: mir-windows mir-windows-386 mir-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

mir:
	$(GORUN) build/ci.go install ./cmd/mir
	@echo "Done building."
	@echo "Run \"$(GOBIN)/mir\" to launch mir."

bootnode:
	$(GORUN) build/ci.go install ./cmd/bootnode
	@echo "Done building."
	@echo "Run \"$(GOBIN)/bootnode\" to launch bootnode."

all:
	$(GORUN) build/ci.go install

# android:
# 	$(GORUN) build/ci.go aar --local
# 	@echo "Done building."
# 	@echo "Import \"$(GOBIN)/mir.aar\" to use the library."
# 	@echo "Import \"$(GOBIN)/mir-sources.jar\" to add javadocs"
# 	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"

# ios:
# 	$(GORUN) build/ci.go xcode --local
# 	@echo "Done building."
# 	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go install golang.org/x/tools/cmd/stringer@latest
	env GOBIN= go install github.com/kevinburke/go-bindata/go-bindata@latest
	env GOBIN= go install github.com/fjl/gencodec@latest
	env GOBIN= go install github.com/golang/protobuf/protoc-gen-go@latest
	env GOBIN= go install ./cmd/abigen
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo) // Mir: not working anymore duew to xgo outdated, see https://github.com/ethereum/go-ethereum/issues/26170

# mir-cross: mir-linux mir-darwin mir-windows mir-android mir-ios
# 	@echo "Full cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-*

# mir-linux: mir-linux-386 mir-linux-amd64 mir-linux-arm mir-linux-mips64 mir-linux-mips64le
# 	@echo "Linux cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-*

# mir-linux-386:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/mir
# 	@echo "Linux 386 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep 386

# mir-linux-amd64:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/mir
# 	@echo "Linux amd64 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep amd64

# mir-linux-arm: mir-linux-arm-5 mir-linux-arm-6 mir-linux-arm-7 mir-linux-arm64
# 	@echo "Linux ARM cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep arm

# mir-linux-arm-5:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/mir
# 	@echo "Linux ARMv5 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep arm-5

# mir-linux-arm-6:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/mir
# 	@echo "Linux ARMv6 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep arm-6

# mir-linux-arm-7:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/mir
# 	@echo "Linux ARMv7 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep arm-7

# mir-linux-arm64:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/mir
# 	@echo "Linux ARM64 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep arm64

# mir-linux-mips:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/mir
# 	@echo "Linux MIPS cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep mips

# mir-linux-mipsle:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/mir
# 	@echo "Linux MIPSle cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep mipsle

# mir-linux-mips64:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/mir
# 	@echo "Linux MIPS64 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep mips64

# mir-linux-mips64le:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/mir
# 	@echo "Linux MIPS64le cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-linux-* | grep mips64le

# mir-darwin: mir-darwin-386 mir-darwin-amd64
# 	@echo "Darwin cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-darwin-*

# mir-darwin-386:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/mir
# 	@echo "Darwin 386 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-darwin-* | grep 386

# mir-darwin-amd64:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/mir
# 	@echo "Darwin amd64 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-darwin-* | grep amd64

# mir-windows: mir-windows-386 mir-windows-amd64
# 	@echo "Windows cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-windows-*

# mir-windows-386:
# 	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/mir
# 	@echo "Windows 386 cross compilation done:"
# 	@ls -ld $(GOBIN)/mir-windows-* | grep 386

mir-windows-amd64:
	$(GORUN) build/ci.go install -os windows -arch amd64 ./cmd/mir
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/mir-windows-* | grep amd64
