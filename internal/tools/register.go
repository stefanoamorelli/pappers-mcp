package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

// RegisterAll registers all 24 Pappers API tools with the MCP server.
func RegisterAll(s *mcp.Server, c client.PappersClient) {
	// 1. Company data
	s.AddTool(companyTool(), companyHandler(c))

	// 2. Association data
	s.AddTool(associationTool(), associationHandler(c))

	// 3-7. Search tools
	s.AddTool(searchCompaniesTool(), searchCompaniesHandler(c))
	s.AddTool(searchDirectorsTool(), searchDirectorsHandler(c))
	s.AddTool(searchBeneficiariesTool(), searchBeneficiariesHandler(c))
	s.AddTool(searchDocumentsTool(), searchDocumentsHandler(c))
	s.AddTool(searchPublicationsTool(), searchPublicationsHandler(c))

	// 8. Suggestions
	s.AddTool(suggestCompaniesTool(), suggestCompaniesHandler(c))

	// 9. Annual accounts
	s.AddTool(annualAccountsTool(), annualAccountsHandler(c))

	// 10. Corporate map
	s.AddTool(corporateMapTool(), corporateMapHandler(c))

	// 11. Compliance
	s.AddTool(pepSanctionsTool(), pepSanctionsHandler(c))

	// 12. Document download (by token)
	s.AddTool(downloadDocumentTool(), downloadDocumentHandler(c))

	// 13-19. Company document downloads
	s.AddTool(pappersExtractTool(), companyDocumentHandler(c, "/document/extrait_pappers"))
	s.AddTool(inpiExtractTool(), companyDocumentHandler(c, "/document/extrait_inpi"))
	s.AddTool(inseeNoticeTool(), companyDocumentHandler(c, "/document/avis_situation_insee"))
	s.AddTool(companyBylawsTool(), companyDocumentHandler(c, "/document/statuts"))
	s.AddTool(uboDeclarationTool(), companyDocumentHandler(c, "/document/declaration_beneficiaires_effectifs"))
	s.AddTool(financialReportTool(), companyDocumentHandler(c, "/document/rapport_financier"))
	s.AddTool(nonFinancialReportTool(), companyDocumentHandler(c, "/document/rapport_non_financier"))

	// 20-23. Surveillance tools
	s.AddTool(addCompanyWatchTool(), addCompanyWatchHandler(c))
	s.AddTool(addDirectorWatchTool(), addDirectorWatchHandler(c))
	s.AddTool(deleteNotificationsTool(), deleteNotificationsHandler(c))
	s.AddTool(addNotificationInfoTool(), addNotificationInfoHandler(c))

	// 24. API credits
	s.AddTool(apiCreditsTool(), apiCreditsHandler(c))
}
