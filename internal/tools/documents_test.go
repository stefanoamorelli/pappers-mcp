package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestDownloadDocumentHandler_RequiresToken(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := downloadDocumentHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when token is missing")
	}
}

func TestDownloadDocumentHandler_HappyPath(t *testing.T) {
	pdfContent := []byte("%PDF-1.4 test content")
	mock := &testutil.MockPappersClient{
		DownloadDocumentFunc: func(ctx context.Context, token string) ([]byte, string, error) {
			if token != "abc123" {
				t.Errorf("expected token=abc123, got %s", token)
			}
			return pdfContent, "application/pdf", nil
		},
	}

	handler := downloadDocumentHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"token": "abc123"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}

	text := result.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(text, "application/pdf") {
		t.Error("expected content_type in response")
	}
	encoded := base64.StdEncoding.EncodeToString(pdfContent)
	if !strings.Contains(text, encoded) {
		t.Error("expected base64 data in response")
	}
}

func TestCompanyDocumentHandler_RequiresSiren(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := companyDocumentHandler(mock, "/document/extrait_pappers")
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when siren is missing")
	}
}

func TestCompanyDocumentHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		DownloadCompanyDocumentFunc: func(ctx context.Context, path string, params url.Values) ([]byte, string, error) {
			if path != "/document/extrait_pappers" {
				t.Errorf("expected path /document/extrait_pappers, got %s", path)
			}
			if params.Get("siren") != "443061841" {
				t.Errorf("expected siren=443061841, got %s", params.Get("siren"))
			}
			return []byte("PDF"), "application/pdf", nil
		},
	}

	handler := companyDocumentHandler(mock, "/document/extrait_pappers")
	result, err := handler(context.Background(), makeReq(map[string]any{"siren": "443061841"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestFormatBinaryResponse(t *testing.T) {
	data := []byte("hello")
	result := formatBinaryResponse(data, "text/plain")

	var parsed map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if parsed["content_type"] != "text/plain" {
		t.Errorf("expected content_type=text/plain, got %v", parsed["content_type"])
	}
	if parsed["size_bytes"] != float64(5) {
		t.Errorf("expected size_bytes=5, got %v", parsed["size_bytes"])
	}
	encoded := base64.StdEncoding.EncodeToString(data)
	if parsed["data_base64"] != encoded {
		t.Error("data_base64 mismatch")
	}
}
