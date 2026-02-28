package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

// listToolNames is a test helper that registers tools with the given filter,
// connects via in-memory transports, and returns the registered tool names.
func listToolNames(t *testing.T, filter ToolFilter) []string {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "pappers-mcp",
		Version: "2.0.0-test",
	}, nil)

	mock := &testutil.MockPappersClient{}
	RegisterAll(server, mock, filter)

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

	var names []string
	for _, tool := range toolResult.Tools {
		names = append(names, tool.Name)
	}
	return names
}

// TestConformance_AllToolsRegistered verifies all 24 tools are registered
// and listed via the MCP protocol.
func TestConformance_AllToolsRegistered(t *testing.T) {
	names := listToolNames(t, nil)

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

	if len(names) != len(expectedTools) {
		t.Errorf("expected %d tools, got %d", len(expectedTools), len(names))
		for _, name := range names {
			t.Logf("  registered: %s", name)
		}
	}

	registeredNames := make(map[string]bool)
	for _, name := range names {
		registeredNames[name] = true
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
	RegisterAll(server, mock, nil)

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

func TestRegisterAll_EnabledFilter(t *testing.T) {
	filter := NewToolFilter("get_company_data,get_api_credits", "")
	names := listToolNames(t, filter)

	if len(names) != 2 {
		t.Fatalf("expected 2 tools, got %d: %v", len(names), names)
	}

	set := make(map[string]bool)
	for _, n := range names {
		set[n] = true
	}
	if !set["get_company_data"] {
		t.Error("expected get_company_data to be registered")
	}
	if !set["get_api_credits"] {
		t.Error("expected get_api_credits to be registered")
	}
}

func TestRegisterAll_DisabledFilter(t *testing.T) {
	filter := NewToolFilter("", "get_api_credits")
	names := listToolNames(t, filter)

	if len(names) != 23 {
		t.Fatalf("expected 23 tools, got %d: %v", len(names), names)
	}

	for _, n := range names {
		if n == "get_api_credits" {
			t.Error("get_api_credits should not be registered when disabled")
		}
	}
}

func TestNewToolFilter(t *testing.T) {
	t.Run("nil when neither set", func(t *testing.T) {
		f := NewToolFilter("", "")
		if f != nil {
			t.Error("expected nil filter when neither enabled nor disabled is set")
		}
	})

	t.Run("enabled only", func(t *testing.T) {
		f := NewToolFilter("a,b", "")
		if f == nil {
			t.Fatal("expected non-nil filter")
		}
		if !f("a") || !f("b") {
			t.Error("expected a and b to pass filter")
		}
		if f("c") {
			t.Error("expected c to be rejected")
		}
	})

	t.Run("disabled only", func(t *testing.T) {
		f := NewToolFilter("", "x")
		if f == nil {
			t.Fatal("expected non-nil filter")
		}
		if f("x") {
			t.Error("expected x to be rejected")
		}
		if !f("y") {
			t.Error("expected y to pass filter")
		}
	})

	t.Run("enabled takes precedence", func(t *testing.T) {
		f := NewToolFilter("a", "b")
		if f == nil {
			t.Fatal("expected non-nil filter")
		}
		if !f("a") {
			t.Error("expected a to pass")
		}
		if f("b") {
			t.Error("expected b to be rejected (not in enabled list)")
		}
	})

	t.Run("whitespace trimming", func(t *testing.T) {
		f := NewToolFilter(" a , b ", "")
		if f == nil {
			t.Fatal("expected non-nil filter")
		}
		if !f("a") || !f("b") {
			t.Error("expected trimmed names to pass")
		}
	})
}
