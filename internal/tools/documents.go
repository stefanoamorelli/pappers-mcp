package tools

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

// --- Download Document (by token) ---

func downloadDocumentTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "download_document",
		Description: "Download a document from Pappers using a download token obtained from search results. Returns the document as base64-encoded data with content type metadata.",
		InputSchema: objectSchema(map[string]any{
			"token": prop("string", "Download token for the document"),
		}, []string{"token"}),
	}
}

func downloadDocumentHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		token := getString(args, "token", "")
		if token == "" {
			return toolError("Parameter 'token' is required"), nil
		}

		data, contentType, err := c.DownloadDocument(ctx, token)
		if err != nil {
			return toolErrorf("Document download failed: %v", err), nil
		}

		return toolText(formatBinaryResponse(data, contentType)), nil
	}
}

// --- Company document helpers ---

// companyDocumentTool creates a tool definition for a company document download endpoint.
func companyDocumentTool(name, description, path string) *mcp.Tool {
	return &mcp.Tool{
		Name:        name,
		Description: description,
		InputSchema: objectSchema(map[string]any{
			"siren": prop("string", "SIREN number (9 digits)"),
		}, []string{"siren"}),
	}
}

// companyDocumentHandler creates a handler for downloading a company document.
func companyDocumentHandler(c client.PappersClient, path string) mcp.ToolHandler {
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

		data, contentType, err := c.DownloadCompanyDocument(ctx, path, params)
		if err != nil {
			return toolErrorf("Document download failed: %v", err), nil
		}

		return toolText(formatBinaryResponse(data, contentType)), nil
	}
}

// --- Individual document tools ---

func pappersExtractTool() *mcp.Tool {
	return companyDocumentTool(
		"get_pappers_extract",
		"Download a Pappers company extract (extrait Pappers) as a PDF. Provides a comprehensive overview of the company including legal info, directors, financials, and more.",
		"/document/extrait_pappers",
	)
}

func inpiExtractTool() *mcp.Tool {
	return companyDocumentTool(
		"get_inpi_extract",
		"Download an INPI extract (extrait INPI / Kbis equivalent) for a company. Official extract from the French National Institute of Industrial Property.",
		"/document/extrait_inpi",
	)
}

func inseeNoticeTool() *mcp.Tool {
	return companyDocumentTool(
		"get_insee_notice",
		"Download the INSEE situation notice (avis de situation INSEE) for a company. Official document from the French National Institute of Statistics.",
		"/document/avis_situation_insee",
	)
}

func companyBylawsTool() *mcp.Tool {
	return companyDocumentTool(
		"get_company_bylaws",
		"Download the company bylaws (statuts) document. Contains the founding articles and rules governing the company.",
		"/document/statuts",
	)
}

func uboDeclarationTool() *mcp.Tool {
	return companyDocumentTool(
		"get_ubo_declaration",
		"Download the declaration of ultimate beneficial owners (déclaration des bénéficiaires effectifs) for a company.",
		"/document/declaration_beneficiaires_effectifs",
	)
}

func financialReportTool() *mcp.Tool {
	return companyDocumentTool(
		"get_financial_report",
		"Download the financial report (rapport financier) for a company.",
		"/document/rapport_financier",
	)
}

func nonFinancialReportTool() *mcp.Tool {
	return companyDocumentTool(
		"get_non_financial_report",
		"Download the non-financial performance report (rapport extra-financier) for a company. Includes ESG and CSR information.",
		"/document/rapport_non_financier",
	)
}

// formatBinaryResponse formats binary data for text-based MCP output.
func formatBinaryResponse(data []byte, contentType string) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf(`{"content_type": %q, "size_bytes": %d, "data_base64": %q}`, contentType, len(data), encoded)
}
