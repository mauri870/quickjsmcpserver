package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/mcptest"
	"github.com/mark3labs/mcp-go/server"
)

func TestServer(t *testing.T) {
	ctx := t.Context()

	testCases := []struct {
		name    string
		code    string
		want    string
		wantErr bool
	}{
		{
			name: "simple op",
			code: "1+1",
			want: "2",
		},
		{
			name:    "function call",
			code:    "const fib = n => n <= 1 ? n : fib(n - 1) + fib(n - 2); fib(10)",
			want:    "55",
			wantErr: false,
		},
		{
			name:    "syntax error",
			code:    "1+",
			want:    "SyntaxError: unexpected token in expression",
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			srv, err := mcptest.NewServer(t, server.ServerTool{
				Tool:    createQuickjsTool(),
				Handler: quickjsHandler,
			})
			if err != nil {
				t.Fatal(err)
			}
			defer srv.Close()

			client := srv.Client()

			var req mcp.CallToolRequest
			req.Params.Name = "quickjs"
			req.Params.Arguments = map[string]any{
				"code": tt.code,
			}

			result, err := client.CallTool(ctx, req)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("CallTool: Got error %q, want %q", err, tt.want)
				} else if !strings.Contains(err.Error(), tt.want) {
					t.Fatalf("CallTool: Got error %q, want %q", err, tt.want)
				}
				return
			}

			got, err := resultToString(result)
			if err != nil {
				t.Fatalf("resultToString: %v", err)
			}

			if got != tt.want {
				t.Fatalf("Got %q, want %q", got, tt.want)
			}
		})
	}
}

func resultToString(result *mcp.CallToolResult) (string, error) {
	var b strings.Builder

	for _, content := range result.Content {
		text, ok := content.(mcp.TextContent)
		if !ok {
			return "", fmt.Errorf("unsupported content type: %T", content)
		}
		b.WriteString(text.Text)
	}

	if result.IsError {
		return "", fmt.Errorf("%s", b.String())
	}

	return b.String(), nil
}
