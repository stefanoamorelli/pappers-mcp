package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

func pepSanctionsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "check_pep_sanctions",
		Description: "Check if a person is a Politically Exposed Person (PEP) or is listed on international sanctions lists. Used for AML/KYC compliance screening.",
		InputSchema: objectSchema(map[string]any{
			"nom":              prop("string", "Last name of the person to check"),
			"prenom":           prop("string", "First name of the person to check"),
			"date_de_naissance": prop("string", "Birth date (YYYY-MM-DD)"),
		}, []string{"nom", "prenom"}),
	}
}

func pepSanctionsHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		nom := getString(args, "nom", "")
		prenom := getString(args, "prenom", "")
		if nom == "" || prenom == "" {
			return toolError("Parameters 'nom' and 'prenom' are required"), nil
		}

		params := url.Values{}
		params.Set("nom", nom)
		params.Set("prenom", prenom)
		setString(params, "date_de_naissance", getString(args, "date_de_naissance", ""))

		data, err := c.CheckPEP(ctx, params)
		if err != nil {
			return toolErrorf("PEP/sanctions check failed: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
