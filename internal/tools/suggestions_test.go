package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestSuggestHandler_RequiresQ(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := suggestCompaniesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when q is missing")
	}
}

func TestSuggestHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		SuggestFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("q") != "goo" {
				t.Errorf("expected q=goo, got %s", params.Get("q"))
			}
			return testutil.SuggestionsFixture(), nil
		},
	}

	handler := suggestCompaniesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"q": "goo"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}
