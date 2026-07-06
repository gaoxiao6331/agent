package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer("Go MCP Demo Server", "1.0.0")

	greetTool := mcp.NewTool("greet",
		mcp.WithDescription("对指定的人进行友好地问候"),
		mcp.WithString("name", mcp.Required(), mcp.Description("要问候的人的名字")),
	)

	s.AddTool(greetTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("无效的参数")
		}
		name, ok := args["name"].(string)
		if !ok {
			return nil, fmt.Errorf("参数 name 必须是字符串类型")
		}
		greeting := fmt.Sprintf("你好，%s！欢迎使用由 Golang 实现的 Model Context Protocol (MCP) 服务！", name)
		return mcp.NewToolResultText(greeting), nil
	})

	addTool := mcp.NewTool("add",
		mcp.WithDescription("计算两个数字的和"),
		mcp.WithNumber("a", mcp.Required(), mcp.Description("第一个加数")),
		mcp.WithNumber("b", mcp.Required(), mcp.Description("第二个加数")),
	)

	s.AddTool(addTool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := req.Params.Arguments.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("无效的参数")
		}
		a, okA := args["a"].(float64)
		b, okB := args["b"].(float64)
		if !okA || !okB {
			return nil, fmt.Errorf("参数 a 和 b 必须是数字类型")
		}
		result := a + b
		return mcp.NewToolResultText(fmt.Sprintf("计算结果: %f + %f = %f", a, b, result)), nil
	})

	port := ":8080"
	fmt.Fprintf(os.Stderr, "正在启动 Go MCP 服务端，监听端口 %s\n", port)

	sse := server.NewSSEServer(s)
	if err := sse.Start(port); err != nil {
		fmt.Fprintf(os.Stderr, "服务运行出错: %v\n", err)
		os.Exit(1)
	}
}
