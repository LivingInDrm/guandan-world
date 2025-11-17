.PHONY: proto proto-go proto-js proto-ts clean-proto help

help:
	@echo "Available targets:"
	@echo "  proto        - Generate all proto files (Go + frontend)"
	@echo "  proto-go     - Generate Go proto files only"
	@echo "  proto-js     - Generate JavaScript proto files"
	@echo "  proto-ts     - Generate TypeScript proto files"
	@echo "  clean-proto  - Clean generated proto files"

proto: proto-go
	@echo "✅ All proto files generated successfully"

proto-go:
	@echo "🔨 Generating Go proto files..."
	@mkdir -p proto/event
	protoc --proto_path=proto \
		--go_out=. \
		--go_opt=paths=source_relative \
		--plugin=protoc-gen-go=$$(go env GOPATH)/bin/protoc-gen-go \
		proto/event.proto
	@if [ -f event.pb.go ]; then mv event.pb.go proto/event/; fi
	@echo "✅ Go proto files generated at proto/event/event.pb.go"

proto-js:
	@echo "🔨 Generating JavaScript proto files..."
	@mkdir -p frontend/src/proto
	protoc --proto_path=proto \
		--js_out=import_style=commonjs,binary:frontend/src/proto \
		proto/event.proto
	@echo "✅ JavaScript proto files generated"

proto-ts:
	@echo "🔨 Generating TypeScript proto files..."
	@mkdir -p frontend/src/proto
	protoc --proto_path=proto \
		--plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
		--ts_out=frontend/src/proto \
		proto/event.proto
	@echo "✅ TypeScript proto files generated"

clean-proto:
	@echo "🧹 Cleaning generated proto files..."
	@rm -f proto/event/*.pb.go
	@rm -f frontend/src/proto/event_pb.*
	@echo "✅ Cleaned proto files"
