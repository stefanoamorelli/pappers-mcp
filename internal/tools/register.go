package tools

import (
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

// ToolFilter decides whether a tool with the given name should be registered.
// Return true to include the tool, false to exclude it.
type ToolFilter func(toolName string) bool

// NewToolFilter builds a ToolFilter from comma-separated allowlist/blocklist
// strings (typically read from environment variables).
//
// Rules:
//   - If enabled is non-empty, only tool names listed in enabled are included.
//   - Otherwise, if disabled is non-empty, all tools except those listed are included.
//   - If both are empty, nil is returned (all tools registered).
//   - enabled takes precedence when both are set.
func NewToolFilter(enabled, disabled string) ToolFilter {
	if enabled != "" {
		set := parseCSV(enabled)
		return func(name string) bool {
			_, ok := set[name]
			return ok
		}
	}
	if disabled != "" {
		set := parseCSV(disabled)
		return func(name string) bool {
			_, ok := set[name]
			return !ok
		}
	}
	return nil
}

// parseCSV splits a comma-separated string into a set, trimming whitespace
// and ignoring empty entries.
func parseCSV(s string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, part := range strings.Split(s, ",") {
		if v := strings.TrimSpace(part); v != "" {
			set[v] = struct{}{}
		}
	}
	return set
}

// toolEntry pairs a tool definition with its handler for registration.
type toolEntry struct {
	tool    *mcp.Tool
	handler mcp.ToolHandler
}

// RegisterAll registers Pappers API tools with the MCP server.
// When filter is non-nil, only tools accepted by the filter are registered.
func RegisterAll(s *mcp.Server, c client.PappersClient, filter ToolFilter) {
	entries := []toolEntry{
		// 1. Company data
		{companyTool(), companyHandler(c)},
		// 2. Association data
		{associationTool(), associationHandler(c)},
		// 3-7. Search tools
		{searchCompaniesTool(), searchCompaniesHandler(c)},
		{searchDirectorsTool(), searchDirectorsHandler(c)},
		{searchBeneficiariesTool(), searchBeneficiariesHandler(c)},
		{searchDocumentsTool(), searchDocumentsHandler(c)},
		{searchPublicationsTool(), searchPublicationsHandler(c)},
		// 8. Suggestions
		{suggestCompaniesTool(), suggestCompaniesHandler(c)},
		// 9. Annual accounts
		{annualAccountsTool(), annualAccountsHandler(c)},
		// 10. Corporate map
		{corporateMapTool(), corporateMapHandler(c)},
		// 11. Compliance
		{pepSanctionsTool(), pepSanctionsHandler(c)},
		// 12. Document download (by token)
		{downloadDocumentTool(), downloadDocumentHandler(c)},
		// 13-19. Company document downloads
		{pappersExtractTool(), companyDocumentHandler(c, "/document/extrait_pappers")},
		{inpiExtractTool(), companyDocumentHandler(c, "/document/extrait_inpi")},
		{inseeNoticeTool(), companyDocumentHandler(c, "/document/avis_situation_insee")},
		{companyBylawsTool(), companyDocumentHandler(c, "/document/statuts")},
		{uboDeclarationTool(), companyDocumentHandler(c, "/document/declaration_beneficiaires_effectifs")},
		{financialReportTool(), companyDocumentHandler(c, "/document/rapport_financier")},
		{nonFinancialReportTool(), companyDocumentHandler(c, "/document/rapport_non_financier")},
		// 20-23. Surveillance tools
		{addCompanyWatchTool(), addCompanyWatchHandler(c)},
		{addDirectorWatchTool(), addDirectorWatchHandler(c)},
		{deleteNotificationsTool(), deleteNotificationsHandler(c)},
		{addNotificationInfoTool(), addNotificationInfoHandler(c)},
		// 24. API credits
		{apiCreditsTool(), apiCreditsHandler(c)},
	}

	for _, e := range entries {
		if filter != nil && !filter(e.tool.Name) {
			continue
		}
		s.AddTool(e.tool, e.handler)
	}
}
