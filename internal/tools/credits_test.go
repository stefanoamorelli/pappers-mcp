package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestAPICreditsHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		GetAPICreditsFunc: func(ctx context.Context) (json.RawMessage, error) {
			return testutil.CreditsFixture(), nil
		},
	}

	handler := apiCreditsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}

	text := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(text, "jetons_utilises") {
		t.Error("expected jetons_utilises in response")
	}
}
