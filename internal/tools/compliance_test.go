package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestPEPHandler_RequiresNomAndPrenom(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := pepSanctionsHandler(mock)

	tests := []struct {
		name string
		args map[string]any
	}{
		{"missing both", map[string]any{}},
		{"missing prenom", map[string]any{"nom": "Doe"}},
		{"missing nom", map[string]any{"prenom": "John"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := handler(context.Background(), makeReq(tt.args))
			if err != nil {
				t.Fatal(err)
			}
			if !result.IsError {
				t.Fatal("expected error")
			}
		})
	}
}

func TestPEPHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		CheckPEPFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("nom") != "Macron" {
				t.Errorf("expected nom=Macron, got %s", params.Get("nom"))
			}
			if params.Get("prenom") != "Emmanuel" {
				t.Errorf("expected prenom=Emmanuel, got %s", params.Get("prenom"))
			}
			return json.RawMessage(`{"pep":true}`), nil
		},
	}

	handler := pepSanctionsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"nom":    "Macron",
		"prenom": "Emmanuel",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}
