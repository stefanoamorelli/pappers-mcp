package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

func corporateMapTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_corporate_map",
		Description: "Retrieve the corporate map (cartographie) for a French company. Shows parent companies, subsidiaries, directors, and their relationships in a graph structure.",
		InputSchema: objectSchema(map[string]any{
			"siren": prop("string", "SIREN number (9 digits) of the company"),
		}, []string{"siren"}),
	}
}

func corporateMapHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		siren := getString(args, "siren", "")
		if siren == "" {
			return toolError("Parameter 'siren' is required"), nil
		}

		params := url.Values{}
		params.Set("siren", siren)

		data, err := c.GetCorporateMap(ctx, params)
		if err != nil {
			return toolErrorf("Failed to retrieve corporate map: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
