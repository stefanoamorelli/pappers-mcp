package tools

import (
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestExtractFirstSiren_Valid(t *testing.T) {
	data := json.RawMessage(`{"resultats":[{"siren":"443061841","nom_entreprise":"GOOGLE FRANCE"}],"total":1}`)
	siren, err := extractFirstSiren(data)
	if err != nil {
		t.Fatal(err)
	}
	if siren != "443061841" {
		t.Errorf("expected 443061841, got %s", siren)
	}
}

func TestExtractFirstSiren_Empty(t *testing.T) {
	data := json.RawMessage(`{"resultats":[],"total":0}`)
	_, err := extractFirstSiren(data)
	if err == nil {
		t.Fatal("expected error for empty results")
	}
}

func TestExtractFirstSiren_NoSiren(t *testing.T) {
	data := json.RawMessage(`{"resultats":[{"nom_entreprise":"TEST"}],"total":1}`)
	_, err := extractFirstSiren(data)
	if err == nil {
		t.Fatal("expected error when siren is empty")
	}
}

func TestExtractFirstSiren_MalformedJSON(t *testing.T) {
	data := json.RawMessage(`not json`)
	_, err := extractFirstSiren(data)
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}

func TestGetString(t *testing.T) {
	args := map[string]any{"name": "test", "count": 42.0}

	if v := getString(args, "name", ""); v != "test" {
		t.Errorf("expected test, got %s", v)
	}
	if v := getString(args, "missing", "default"); v != "default" {
		t.Errorf("expected default, got %s", v)
	}
	if v := getString(args, "count", "fallback"); v != "fallback" {
		t.Errorf("expected fallback for non-string, got %s", v)
	}
}

func TestGetInt(t *testing.T) {
	args := map[string]any{"page": 5.0, "name": "test"}

	v, ok := getInt(args, "page")
	if !ok || v != 5 {
		t.Errorf("expected 5, got %d (ok=%v)", v, ok)
	}

	_, ok = getInt(args, "missing")
	if ok {
		t.Error("expected ok=false for missing key")
	}

	_, ok = getInt(args, "name")
	if ok {
		t.Error("expected ok=false for non-numeric value")
	}
}

func TestGetBool(t *testing.T) {
	args := map[string]any{"active": true, "name": "test"}

	v, ok := getBool(args, "active")
	if !ok || !v {
		t.Errorf("expected true, got %v (ok=%v)", v, ok)
	}

	_, ok = getBool(args, "missing")
	if ok {
		t.Error("expected ok=false for missing key")
	}
}

func TestBuildSearchParams(t *testing.T) {
	args := map[string]any{
		"q":           "google",
		"page":        2.0,
		"par_page":    20.0,
		"code_postal": "75009",
	}

	params := buildSearchParams(args)

	if params.Get("q") != "google" {
		t.Errorf("expected q=google, got %s", params.Get("q"))
	}
	if params.Get("page") != "2" {
		t.Errorf("expected page=2, got %s", params.Get("page"))
	}
	if params.Get("par_page") != "20" {
		t.Errorf("expected par_page=20, got %s", params.Get("par_page"))
	}
	if params.Get("code_postal") != "75009" {
		t.Errorf("expected code_postal=75009, got %s", params.Get("code_postal"))
	}
	// Empty/missing params should not be set
	if params.Get("departement") != "" {
		t.Errorf("expected empty departement, got %s", params.Get("departement"))
	}
}

func TestMergeProperties(t *testing.T) {
	a := map[string]any{"x": 1, "y": 2}
	b := map[string]any{"y": 3, "z": 4}
	result := mergeProperties(a, b)

	if result["x"] != 1 {
		t.Errorf("expected x=1, got %v", result["x"])
	}
	if result["y"] != 3 {
		t.Errorf("expected y=3 (overridden), got %v", result["y"])
	}
	if result["z"] != 4 {
		t.Errorf("expected z=4, got %v", result["z"])
	}
}

func TestToolText(t *testing.T) {
	result := toolText("hello")
	if result.IsError {
		t.Error("expected IsError=false")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *TextContent, got %T", result.Content[0])
	}
	if tc.Text != "hello" {
		t.Errorf("expected hello, got %s", tc.Text)
	}
}

func TestToolError(t *testing.T) {
	result := toolError("bad input")
	if !result.IsError {
		t.Error("expected IsError=true")
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected *TextContent, got %T", result.Content[0])
	}
	if tc.Text != "bad input" {
		t.Errorf("expected 'bad input', got %s", tc.Text)
	}
}
