package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestAssociationHandler_BySiren(t *testing.T) {
	mock := &testutil.MockPappersClient{
		GetAssociationFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("siren") != "123456789" {
				t.Errorf("expected siren=123456789, got %s", params.Get("siren"))
			}
			return testutil.AssociationFixture(), nil
		},
	}

	handler := associationHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"siren": "123456789"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
	text := result.Content[0].(*mcp.TextContent).Text
	var resp map[string]any
	if err := json.Unmarshal([]byte(text), &resp); err != nil {
		t.Fatal(err)
	}
	if resp["titre"] != "ASSOCIATION TEST" {
		t.Errorf("expected ASSOCIATION TEST, got %v", resp["titre"])
	}
}

func TestAssociationHandler_NoParams(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := associationHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when no params")
	}
}
