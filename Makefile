
all: build

build:
	go build .

inspector: build
	npx @modelcontextprotocol/inspector ./quickjsmcpserver
