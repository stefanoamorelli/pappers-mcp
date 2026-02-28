package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

func suggestCompaniesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "suggest_companies",
		Description: "Get company name suggestions (autocomplete) from the Pappers API. Useful for interactive search-as-you-type interfaces. Returns a list of matching companies with basic info.",
		InputSchema: objectSchema(map[string]any{
			"q":        prop("string", "Partial company name to autocomplete"),
			"longueur": prop("integer", "Maximum number of suggestions to return (default 10)"),
			"cibles":   prop("string", "Target types: nom_entreprise, denomination, nom_complet, representant, siren, siret (comma-separated)"),
		}, []string{"q"}),
	}
}

func suggestCompaniesHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		q := getString(args, "q", "")
		if q == "" {
			return toolError("Parameter 'q' is required"), nil
		}

		params := url.Values{}
		params.Set("q", q)
		setInt(params, "longueur", args)
		setString(params, "cibles", getString(args, "cibles", ""))

		data, err := c.Suggest(ctx, params)
		if err != nil {
			return toolErrorf("Suggestion request failed: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
