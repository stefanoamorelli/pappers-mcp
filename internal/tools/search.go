package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

// --- Search Companies ---

func searchCompaniesTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "search_companies",
		Description: "Search for French companies in the Pappers database. Supports extensive filtering by location, legal form, financials, activity codes, and more. Returns paginated results.",
		InputSchema: objectSchema(commonSearchProperties(), []string{"q"}),
	}
}

func searchCompaniesHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := buildSearchParams(args)
		if params.Get("q") == "" {
			return toolError("Parameter 'q' (search query) is required"), nil
		}

		data, err := c.SearchCompanies(ctx, params)
		if err != nil {
			return toolErrorf("Search failed: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}

// --- Search Directors ---

func searchDirectorsTool() *mcp.Tool {
	props := mergeProperties(commonSearchProperties(), directorSearchProperties())
	return &mcp.Tool{
		Name:        "search_directors",
		Description: "Search for company directors (dirigeants) in the Pappers database. Filter by name, age, nationality, role, and company attributes. Returns paginated results.",
		InputSchema: objectSchema(props, []string{"q"}),
	}
}

func searchDirectorsHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := buildSearchParams(args)
		if params.Get("q") == "" {
			return toolError("Parameter 'q' (search query) is required"), nil
		}

		// Add director-specific params
		setInt(params, "age_dirigeant_min", args)
		setInt(params, "age_dirigeant_max", args)
		setString(params, "date_de_naissance_dirigeant_min", getString(args, "date_de_naissance_dirigeant_min", ""))
		setString(params, "date_de_naissance_dirigeant_max", getString(args, "date_de_naissance_dirigeant_max", ""))
		setString(params, "nationalite_dirigeant", getString(args, "nationalite_dirigeant", ""))
		setString(params, "qualite_dirigeant", getString(args, "qualite_dirigeant", ""))

		data, err := c.SearchDirectors(ctx, params)
		if err != nil {
			return toolErrorf("Search failed: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}

// --- Search Beneficiaries ---

func searchBeneficiariesTool() *mcp.Tool {
	props := mergeProperties(commonSearchProperties(), beneficiarySearchProperties())
	return &mcp.Tool{
		Name:        "search_beneficiaries",
		Description: "Search for ultimate beneficial owners (UBOs / bénéficiaires effectifs) in the Pappers database. Filter by name, age, nationality, and company attributes. Returns paginated results.",
		InputSchema: objectSchema(props, []string{"q"}),
	}
}

func searchBeneficiariesHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := buildSearchParams(args)
		if params.Get("q") == "" {
			return toolError("Parameter 'q' (search query) is required"), nil
		}

		// Add beneficiary-specific params
		setInt(params, "age_beneficiaire_min", args)
		setInt(params, "age_beneficiaire_max", args)
		setString(params, "date_de_naissance_beneficiaire_min", getString(args, "date_de_naissance_beneficiaire_min", ""))
		setString(params, "date_de_naissance_beneficiaire_max", getString(args, "date_de_naissance_beneficiaire_max", ""))
		setString(params, "nationalite_beneficiaire", getString(args, "nationalite_beneficiaire", ""))
		setString(params, "type_beneficiaire", getString(args, "type_beneficiaire", ""))

		data, err := c.SearchBeneficiaries(ctx, params)
		if err != nil {
			return toolErrorf("Search failed: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}

// --- Search Documents ---

func searchDocumentsTool() *mcp.Tool {
	props := mergeProperties(commonSearchProperties(), map[string]any{
		"type_publication": prop("string", "Document type filter"),
		"date_depot_minimum": prop("string", "Minimum filing date (YYYY-MM-DD)"),
		"date_depot_maximum": prop("string", "Maximum filing date (YYYY-MM-DD)"),
	})
	return &mcp.Tool{
		Name:        "search_documents",
		Description: "Search for company documents (annual accounts, statutes, etc.) in the Pappers database. Filter by company attributes and document type. Returns paginated results with download tokens.",
		InputSchema: objectSchema(props, []string{"q"}),
	}
}

func searchDocumentsHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := buildSearchParams(args)
		if params.Get("q") == "" {
			return toolError("Parameter 'q' (search query) is required"), nil
		}

		setString(params, "type_publication", getString(args, "type_publication", ""))
		setString(params, "date_depot_minimum", getString(args, "date_depot_minimum", ""))
		setString(params, "date_depot_maximum", getString(args, "date_depot_maximum", ""))

		data, err := c.SearchDocuments(ctx, params)
		if err != nil {
			return toolErrorf("Search failed: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}

// --- Search Publications ---

func searchPublicationsTool() *mcp.Tool {
	props := mergeProperties(commonSearchProperties(), map[string]any{
		"type_publication": prop("string", "Publication type filter (e.g. creation, modification, radiation)"),
		"date_publication_minimum": prop("string", "Minimum publication date (YYYY-MM-DD)"),
		"date_publication_maximum": prop("string", "Maximum publication date (YYYY-MM-DD)"),
	})
	return &mcp.Tool{
		Name:        "search_publications",
		Description: "Search for BODACC publications (official company announcements) in the Pappers database. Filter by company attributes and publication type/date. Returns paginated results.",
		InputSchema: objectSchema(props, []string{"q"}),
	}
}

func searchPublicationsHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := buildSearchParams(args)
		if params.Get("q") == "" {
			return toolError("Parameter 'q' (search query) is required"), nil
		}

		setString(params, "type_publication", getString(args, "type_publication", ""))
		setString(params, "date_publication_minimum", getString(args, "date_publication_minimum", ""))
		setString(params, "date_publication_maximum", getString(args, "date_publication_maximum", ""))

		data, err := c.SearchPublications(ctx, params)
		if err != nil {
			return toolErrorf("Search failed: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
