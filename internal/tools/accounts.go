package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

func annualAccountsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_annual_accounts",
		Description: "Retrieve the annual accounts (comptes annuels) for a French company. Returns detailed financial data including balance sheet, income statement, and key ratios.",
		InputSchema: objectSchema(map[string]any{
			"siren":  prop("string", "SIREN number (9 digits)"),
			"siret":  prop("string", "SIRET number (14 digits)"),
			"annee":  prop("integer", "Specific year to retrieve accounts for"),
		}, nil),
	}
}

func annualAccountsHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		siren := getString(args, "siren", "")
		siret := getString(args, "siret", "")

		if siren == "" && siret == "" {
			return toolError("At least one parameter is required: siren or siret"), nil
		}

		params := url.Values{}
		setString(params, "siren", siren)
		setString(params, "siret", siret)
		setInt(params, "annee", args)

		data, err := c.GetAnnualAccounts(ctx, params)
		if err != nil {
			return toolErrorf("Failed to retrieve annual accounts: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
