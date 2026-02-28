package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestAnnualAccountsHandler_RequiresSirenOrSiret(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := annualAccountsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when no siren/siret")
	}
}

func TestAnnualAccountsHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		GetAnnualAccountsFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("siren") != "443061841" {
				t.Errorf("expected siren=443061841, got %s", params.Get("siren"))
			}
			return json.RawMessage(`{"comptes":[]}`), nil
		},
	}

	handler := annualAccountsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"siren": "443061841"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}
