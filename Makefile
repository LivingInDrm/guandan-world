.PHONY: proto proto-check proto-clean proto-deps proto-common proto-game proto-messages

# 检测 protoc-gen-go 是否安装
PROTOC_GEN_GO := $(shell command -v protoc-gen-go 2> /dev/null)
GO_BIN := $(shell go env GOPATH)/bin

# 安装 proto 编译依赖
proto-deps:
	@echo "Installing protoc-gen-go..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@echo ""
	@echo "✅ Dependencies installed successfully!"
	@echo "📍 protoc-gen-go installed at: $(GO_BIN)/protoc-gen-go"
	@echo ""
	@echo "⚠️  If 'make proto' fails, please add the following to your shell profile:"
	@echo "   export PATH=\"\$$PATH:$(GO_BIN)\""
	@echo ""

# 编译所有 proto 文件
proto:
ifndef PROTOC_GEN_GO
	@echo "❌ Error: protoc-gen-go not found in PATH"
	@echo ""
	@echo "Please run one of the following:"
	@echo "  1. Install dependencies: make proto-deps"
	@echo "  2. Add Go bin to PATH:   export PATH=\"\$$PATH:$(GO_BIN)\""
	@echo ""
	@exit 1
endif
	@echo "Compiling proto files..."
	@protoc --proto_path=proto \
		--go_out=proto/gen/go \
		--go_opt=paths=source_relative \
		proto/common/*.proto
	@if [ -d "proto/game" ] && [ -n "$$(ls -A proto/game/*.proto 2>/dev/null)" ]; then \
		protoc --proto_path=proto \
			--go_out=proto/gen/go \
			--go_opt=paths=source_relative \
			proto/game/*.proto; \
	fi
	@if [ -d "proto/messages" ] && [ -n "$$(ls -A proto/messages/*.proto 2>/dev/null)" ]; then \
		protoc --proto_path=proto \
			--go_out=proto/gen/go \
			--go_opt=paths=source_relative \
			proto/messages/*.proto; \
	fi
	@echo "✅ Proto compilation completed."

# 仅编译 common proto
proto-common:
ifndef PROTOC_GEN_GO
	@echo "❌ Error: protoc-gen-go not found in PATH"
	@echo "Run: make proto-deps"
	@exit 1
endif
	@echo "Compiling common proto files..."
	@protoc --proto_path=proto \
		--go_out=proto/gen/go \
		--go_opt=paths=source_relative \
		proto/common/*.proto
	@echo "✅ Common proto compiled."

# 仅编译 game proto (待实现)
proto-game:
ifndef PROTOC_GEN_GO
	@echo "❌ Error: protoc-gen-go not found in PATH"
	@echo "Run: make proto-deps"
	@exit 1
endif
	@echo "Compiling game proto files..."
	@if [ -d "proto/game" ] && [ -n "$$(ls -A proto/game/*.proto 2>/dev/null)" ]; then \
		protoc --proto_path=proto \
			--go_out=proto/gen/go \
			--go_opt=paths=source_relative \
			proto/game/*.proto; \
		echo "✅ Game proto compiled."; \
	else \
		echo "⚠️  No game proto files found, skipping..."; \
	fi

# 仅编译 messages proto (待实现)
proto-messages:
ifndef PROTOC_GEN_GO
	@echo "❌ Error: protoc-gen-go not found in PATH"
	@echo "Run: make proto-deps"
	@exit 1
endif
	@echo "Compiling messages proto files..."
	@if [ -d "proto/messages" ] && [ -n "$$(ls -A proto/messages/*.proto 2>/dev/null)" ]; then \
		protoc --proto_path=proto \
			--go_out=proto/gen/go \
			--go_opt=paths=source_relative \
			proto/messages/*.proto; \
		echo "✅ Messages proto compiled."; \
	else \
		echo "⚠️  No messages proto files found, skipping..."; \
	fi

# 仅检查 proto 语法
proto-check:
	@echo "Checking proto syntax..."
	@protoc --proto_path=proto \
		--descriptor_set_out=/dev/null \
		proto/common/*.proto
	@echo "✅ Proto syntax check passed."

# 清理生成的文件
proto-clean:
	@echo "Cleaning generated proto files..."
	@rm -rf proto/gen/go/common/*.pb.go
	@rm -rf proto/gen/go/game/*.pb.go
	@rm -rf proto/gen/go/messages/*.pb.go
	@echo "✅ Proto files cleaned."
