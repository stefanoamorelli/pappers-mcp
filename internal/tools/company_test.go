package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func makeReq(args map[string]any) *mcp.CallToolRequest {
	data, _ := json.Marshal(args)
	return &mcp.CallToolRequest{
		Params: &mcp.CallToolParamsRaw{
			Arguments: json.RawMessage(data),
		},
	}
}

func TestCompanyHandler_BySiren(t *testing.T) {
	mock := &testutil.MockPappersClient{
		GetCompanyFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("siren") != "443061841" {
				t.Errorf("expected siren=443061841, got %s", params.Get("siren"))
			}
			return testutil.CompanyFixture(), nil
		},
	}

	handler := companyHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"siren": "443061841"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}

	text := result.Content[0].(*mcp.TextContent).Text
	if text == "" {
		t.Error("expected non-empty response")
	}
}

func TestCompanyHandler_BySiret(t *testing.T) {
	mock := &testutil.MockPappersClient{
		GetCompanyFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("siret") != "44306184100047" {
				t.Errorf("expected siret=44306184100047, got %s", params.Get("siret"))
			}
			return testutil.CompanyFixture(), nil
		},
	}

	handler := companyHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"siret": "44306184100047"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestCompanyHandler_ByName(t *testing.T) {
	mock := &testutil.MockPappersClient{
		SearchCompaniesFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("q") != "GOOGLE FRANCE" {
				t.Errorf("expected q=GOOGLE FRANCE, got %s", params.Get("q"))
			}
			return testutil.SearchResultFixture(), nil
		},
		GetCompanyFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("siren") != "443061841" {
				t.Errorf("expected siren from search result, got %s", params.Get("siren"))
			}
			return testutil.CompanyFixture(), nil
		},
	}

	handler := companyHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"company_name": "GOOGLE FRANCE"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestCompanyHandler_NoParams(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := companyHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when no params provided")
	}
}

func TestCompanyHandler_SearchFails(t *testing.T) {
	mock := &testutil.MockPappersClient{
		SearchCompaniesFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	handler := companyHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"company_name": "NONEXISTENT"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when search fails")
	}
}

func TestCompanyHandler_SearchNoResults(t *testing.T) {
	mock := &testutil.MockPappersClient{
		SearchCompaniesFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			return testutil.EmptySearchFixture(), nil
		},
	}

	handler := companyHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"company_name": "NONEXISTENT"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when no results found")
	}
}
