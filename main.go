package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"modernc.org/quickjs"
)

func main() {
	s := server.NewMCPServer(
		"QuickJS MCP Server",
		"0.0.1",
		server.WithToolCapabilities(false),
	)

	s.AddTool(createQuickjsTool(), quickjsHandler)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func createQuickjsTool() mcp.Tool {
	return mcp.NewTool("quickjs",
		mcp.WithDescription("ES2023 JavaScript interpreter"),
		mcp.WithString("code",
			mcp.Required(),
			mcp.Description("Javascript code to execute"),
		),
	)
}

func quickjsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	code, ok := request.Params.Arguments["code"].(string)
	if !ok {
		return nil, errors.New("code must be a string")
	}

	vm, err := quickjs.NewVM()
	if err != nil {
		return nil, err
	}

	r, err := vm.Eval(code, quickjs.EvalGlobal)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprint(r)), nil
}
