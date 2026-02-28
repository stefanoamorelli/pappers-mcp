package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

func associationTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_association",
		Description: "Retrieve data for a French association (non-profit) from the Pappers API. Provide the RNA number or SIREN number.",
		InputSchema: objectSchema(map[string]any{
			"siren":      prop("string", "SIREN number of the association (9 digits)"),
			"numero_rna": prop("string", "RNA number of the association (W + 9 digits)"),
		}, nil),
	}
}

func associationHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		siren := getString(args, "siren", "")
		rna := getString(args, "numero_rna", "")

		if siren == "" && rna == "" {
			return toolError("At least one parameter is required: siren or numero_rna"), nil
		}

		params := url.Values{}
		setString(params, "siren", siren)
		setString(params, "numero_rna", rna)

		data, err := c.GetAssociation(ctx, params)
		if err != nil {
			return toolErrorf("Failed to retrieve association data: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
