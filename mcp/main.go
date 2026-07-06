package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {

	basePath := normalizePath(getEnv("MCP_BASE_PATH", "/mcp"))
	port := getEnv("MCP_PORT", "8080")

	s := server.NewMCPServer("Go MCP Demo Server", "1.0.0")

	// ========== tools ==========
	s.AddTool(
		mcp.NewTool("greet",
			mcp.WithDescription("问候"),
			mcp.WithString("name", mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name := req.Params.Arguments.(map[string]any)["name"].(string)
			return mcp.NewToolResultText("你好，" + name), nil
		},
	)

	s.AddTool(
		mcp.NewTool("add",
			mcp.WithDescription("加法"),
			mcp.WithNumber("a", mcp.Required()),
			mcp.WithNumber("b", mcp.Required()),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.Params.Arguments.(map[string]any)
			a := args["a"].(float64)
			b := args["b"].(float64)
			return mcp.NewToolResultText(fmt.Sprintf("%v", a+b)), nil
		},
	)

	streamSrv := server.NewStreamableHTTPServer(s)

	mux := http.NewServeMux()

	// ✅ 关键修复：同时兼容 /mcp 和 /mcp/
	mux.Handle(basePath, mcpHandler(streamSrv))
	mux.Handle(basePath+"/", mcpHandler(streamSrv))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	fmt.Println("🚀 MCP running")
	fmt.Println("   port:", port)
	fmt.Println("   basePath:", basePath)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// =========================
// MCP handler middleware
// =========================
func mcpHandler(streamSrv *server.StreamableHTTPServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		defer func() {
			if err := recover(); err != nil {
				fmt.Println("❌ panic:", err)
			}
		}()

		fmt.Println("━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("➡️", r.Method, r.URL.Path)

		// ⚠️ 只建议 debug 时打开
		body, _ := io.ReadAll(r.Body)
		if len(body) > 0 {
			fmt.Println("📦 body:", string(body))
		}

		// restore body（避免 MCP 读不到）
		r.Body = io.NopCloser(strings.NewReader(string(body)))

		streamSrv.ServeHTTP(w, r)

		fmt.Println("⬅️ done cost:", time.Since(start))
	})
}

// =========================
// utils
// =========================
func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// 保证 MCP path 不带尾斜杠
func normalizePath(p string) string {
	if p == "" {
		return "/mcp"
	}
	return "/" + strings.Trim(p, "/")
}