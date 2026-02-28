package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

// TestConformance_AllToolsRegistered verifies all 24 tools are registered
// and listed via the MCP protocol.
func TestConformance_AllToolsRegistered(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "pappers-mcp",
		Version: "2.0.0-test",
	}, nil)

	mock := &testutil.MockPappersClient{}
	RegisterAll(server, mock)

	// Connect via in-memory transports
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	// Connect server
	_, err := server.Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect failed: %v", err)
	}

	// Create client and connect
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect failed: %v", err)
	}
	defer session.Close()

	// List tools
	toolResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	expectedTools := []string{
		"get_company_data",
		"get_association",
		"search_companies",
		"search_directors",
		"search_beneficiaries",
		"search_documents",
		"search_publications",
		"suggest_companies",
		"get_annual_accounts",
		"get_corporate_map",
		"check_pep_sanctions",
		"download_document",
		"get_pappers_extract",
		"get_inpi_extract",
		"get_insee_notice",
		"get_company_bylaws",
		"get_ubo_declaration",
		"get_financial_report",
		"get_non_financial_report",
		"add_company_watch",
		"add_director_watch",
		"delete_notifications",
		"add_notification_info",
		"get_api_credits",
	}

	if len(toolResult.Tools) != len(expectedTools) {
		t.Errorf("expected %d tools, got %d", len(expectedTools), len(toolResult.Tools))
		for _, tool := range toolResult.Tools {
			t.Logf("  registered: %s", tool.Name)
		}
	}

	registeredNames := make(map[string]bool)
	for _, tool := range toolResult.Tools {
		registeredNames[tool.Name] = true
	}

	for _, name := range expectedTools {
		if !registeredNames[name] {
			t.Errorf("expected tool %q to be registered", name)
		}
	}
}

// TestConformance_ToolsHaveDescriptions verifies all tools have descriptions.
func TestConformance_ToolsHaveDescriptions(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "pappers-mcp",
		Version: "2.0.0-test",
	}, nil)

	mock := &testutil.MockPappersClient{}
	RegisterAll(server, mock)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	_, err := server.Connect(context.Background(), serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect failed: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(context.Background(), clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect failed: %v", err)
	}
	defer session.Close()

	toolResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	for _, tool := range toolResult.Tools {
		if tool.Description == "" {
			t.Errorf("tool %q has no description", tool.Name)
		}
		if tool.InputSchema == nil {
			t.Errorf("tool %q has no input schema", tool.Name)
		}
	}
}
