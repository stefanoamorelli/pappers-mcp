package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

func apiCreditsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_api_credits",
		Description: "Check your Pappers API token usage and remaining credits. Returns information about your API plan, usage statistics, and credit balance.",
		InputSchema: objectSchema(map[string]any{}, nil),
	}
}

func apiCreditsHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		data, err := c.GetAPICredits(ctx)
		if err != nil {
			return toolErrorf("Failed to retrieve API credits: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
