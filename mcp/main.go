package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// 创建一个新的 MCP 服务端，指定服务名称和版本号
	s := server.NewMCPServer("Go MCP Demo Server", "1.0.0")

	// 1. 注册一个 greet（打招呼）工具
	greetTool := mcp.NewTool("greet",
		mcp.WithDescription("对指定的人进行友好地问候"),
		mcp.WithString("name", mcp.Required(), mcp.WithDescription("要问候的人的名字")),
	)

	s.AddTool(greetTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := req.Arguments["name"].(string)
		if !ok {
			return nil, fmt.Errorf("参数 name 必须是字符串类型")
		}
		greeting := fmt.Sprintf("你好，%s！欢迎使用由 Golang 实现的 Model Context Protocol (MCP) 服务！", name)
		return mcp.NewToolResultText(greeting), nil
	})

	// 2. 注册一个简单的 add（加法计算）工具
	addTool := mcp.NewTool("add",
		mcp.WithDescription("计算两个数字的和"),
		mcp.WithNumber("a", mcp.Required(), mcp.WithDescription("第一个加数")),
		mcp.WithNumber("b", mcp.Required(), mcp.WithDescription("第二个加数")),
	)

	s.AddTool(addTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		a, okA := req.Arguments["a"].(float64)
		b, okB := req.Arguments["b"].(float64)
		if !okA || !okB {
			return nil, fmt.Errorf("参数 a 和 b 必须是数字类型")
		}
		result := a + b
		return mcp.NewToolResultText(fmt.Sprintf("计算结果: %f + %f = %f", a, b, result)), nil
	})

	// 启动 Stdio 服务，供 MCP Host (如 Claude Desktop) 调起
	fmt.Fprintln(os.Stderr, "正在启动 Go MCP 服务端...")
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "服务运行出错: %v", err)
		os.Exit(1)
	}
}
