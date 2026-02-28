package tools

import (
	"context"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

func companyTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "get_company_data",
		Description: "Retrieve complete data for a French company from the Pappers API. Provide at least one search parameter: SIREN, SIRET, or company_name. When company_name is provided, a search is performed first, then the full data for the top result is returned.",
		InputSchema: objectSchema(map[string]any{
			"siren":        prop("string", "SIREN number (9 digits)"),
			"siret":        prop("string", "SIRET number (14 digits)"),
			"company_name": prop("string", "Company name or trade name"),
		}, nil),
	}
}

func companyHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		siren := getString(args, "siren", "")
		siret := getString(args, "siret", "")
		companyName := getString(args, "company_name", "")

		if siren == "" && siret == "" && companyName == "" {
			return toolError("At least one parameter is required: siren, siret, or company_name"), nil
		}

		// If company_name is provided (and no SIREN/SIRET), search first
		if siren == "" && siret == "" && companyName != "" {
			searchParams := url.Values{}
			searchParams.Set("q", companyName)
			searchParams.Set("par_page", "1")
			searchData, err := c.SearchCompanies(ctx, searchParams)
			if err != nil {
				return toolErrorf("Search failed: %v", err), nil
			}
			foundSiren, err := extractFirstSiren(searchData)
			if err != nil {
				return toolErrorf("No company found for name '%s': %v", companyName, err), nil
			}
			siren = foundSiren
		}

		params := url.Values{}
		if siren != "" {
			params.Set("siren", siren)
		}
		if siret != "" {
			params.Set("siret", siret)
		}

		data, err := c.GetCompany(ctx, params)
		if err != nil {
			return toolErrorf("Failed to retrieve company data: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
