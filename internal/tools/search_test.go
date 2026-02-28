package tools

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestSearchCompaniesHandler_RequiresQ(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := searchCompaniesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when q is missing")
	}
}

func TestSearchCompaniesHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		SearchCompaniesFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("q") != "google" {
				t.Errorf("expected q=google, got %s", params.Get("q"))
			}
			if params.Get("code_postal") != "75009" {
				t.Errorf("expected code_postal=75009, got %s", params.Get("code_postal"))
			}
			return testutil.SearchResultFixture(), nil
		},
	}

	handler := searchCompaniesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"q":           "google",
		"code_postal": "75009",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestSearchDirectorsHandler_RequiresQ(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := searchDirectorsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when q is missing")
	}
}

func TestSearchDirectorsHandler_WithDirectorParams(t *testing.T) {
	mock := &testutil.MockPappersClient{
		SearchDirectorsFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("q") != "dupont" {
				t.Errorf("expected q=dupont, got %s", params.Get("q"))
			}
			if params.Get("nationalite_dirigeant") != "Française" {
				t.Errorf("expected nationalite_dirigeant=Française, got %s", params.Get("nationalite_dirigeant"))
			}
			return json.RawMessage(`{"resultats":[]}`), nil
		},
	}

	handler := searchDirectorsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"q":                      "dupont",
		"nationalite_dirigeant": "Française",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestSearchBeneficiariesHandler_RequiresQ(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := searchBeneficiariesHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when q is missing")
	}
}

func TestSearchDocumentsHandler_RequiresQ(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := searchDocumentsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when q is missing")
	}
}

func TestSearchPublicationsHandler_RequiresQ(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := searchPublicationsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when q is missing")
	}
}

func TestSearchPublicationsHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		SearchPublicationsFunc: func(ctx context.Context, params url.Values) (json.RawMessage, error) {
			if params.Get("q") != "test" {
				t.Errorf("expected q=test, got %s", params.Get("q"))
			}
			if params.Get("type_publication") != "creation" {
				t.Errorf("expected type_publication=creation, got %s", params.Get("type_publication"))
			}
			return json.RawMessage(`{"resultats":[]}`), nil
		},
	}

	handler := searchPublicationsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"q":                "test",
		"type_publication": "creation",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}
