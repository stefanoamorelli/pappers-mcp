// Package tools provides MCP tool definitions and handlers for the Pappers API.
package tools

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// --- Result helpers ---

// toolText returns a successful CallToolResult with text content.
func toolText(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: text},
		},
	}
}

// toolError returns a CallToolResult representing a tool error.
func toolError(msg string) *mcp.CallToolResult {
	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
		IsError: true,
	}
	return result
}

// toolErrorf returns a formatted tool error.
func toolErrorf(format string, args ...any) *mcp.CallToolResult {
	return toolError(fmt.Sprintf(format, args...))
}

// --- Argument extraction ---

// extractArgs unmarshals the raw arguments from a CallToolRequest into a map.
func extractArgs(req *mcp.CallToolRequest) (map[string]any, error) {
	if len(req.Params.Arguments) == 0 {
		return map[string]any{}, nil
	}
	var args map[string]any
	if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	return args, nil
}

// getString extracts a string value from args, returning defaultVal if missing.
func getString(args map[string]any, key, defaultVal string) string {
	v, ok := args[key]
	if !ok {
		return defaultVal
	}
	s, ok := v.(string)
	if !ok {
		return defaultVal
	}
	return s
}

// getInt extracts an integer value from args.
func getInt(args map[string]any, key string) (int, bool) {
	v, ok := args[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int(n), true
	case json.Number:
		i, err := n.Int64()
		if err != nil {
			return 0, false
		}
		return int(i), true
	}
	return 0, false
}

// getFloat extracts a float64 value from args.
func getFloat(args map[string]any, key string) (float64, bool) {
	v, ok := args[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, false
		}
		return f, true
	}
	return 0, false
}

// getBool extracts a bool value from args.
func getBool(args map[string]any, key string) (bool, bool) {
	v, ok := args[key]
	if !ok {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

// --- Parameter building ---

// setString sets a URL parameter if the value is non-empty.
func setString(params url.Values, key, value string) {
	if value != "" {
		params.Set(key, value)
	}
}

// setInt sets a URL parameter from an int if present in args.
func setInt(params url.Values, key string, args map[string]any) {
	if v, ok := getInt(args, key); ok {
		params.Set(key, strconv.Itoa(v))
	}
}

// setFloat sets a URL parameter from a float if present in args.
func setFloat(params url.Values, key string, args map[string]any) {
	if v, ok := getFloat(args, key); ok {
		params.Set(key, strconv.FormatFloat(v, 'f', -1, 64))
	}
}

// setBool sets a URL parameter from a bool if present in args.
func setBool(params url.Values, key string, args map[string]any) {
	if v, ok := getBool(args, key); ok {
		params.Set(key, strconv.FormatBool(v))
	}
}

// buildSearchParams extracts common search parameters from args into url.Values.
func buildSearchParams(args map[string]any) url.Values {
	params := url.Values{}
	setString(params, "q", getString(args, "q", ""))
	setInt(params, "page", args)
	setInt(params, "par_page", args)
	setString(params, "precision", getString(args, "precision", ""))
	setString(params, "code_naf", getString(args, "code_naf", ""))
	setString(params, "departement", getString(args, "departement", ""))
	setString(params, "region", getString(args, "region", ""))
	setString(params, "code_postal", getString(args, "code_postal", ""))
	setString(params, "convention_collective", getString(args, "convention_collective", ""))
	setString(params, "categorie_juridique", getString(args, "categorie_juridique", ""))
	setBool(params, "entreprise_cessee", args)
	setString(params, "statut_rcs", getString(args, "statut_rcs", ""))
	setString(params, "objet_social", getString(args, "objet_social", ""))
	setString(params, "date_creation_minimum", getString(args, "date_creation_minimum", ""))
	setString(params, "date_creation_maximum", getString(args, "date_creation_maximum", ""))
	setString(params, "tranche_effectif", getString(args, "tranche_effectif", ""))
	setString(params, "type_entreprise", getString(args, "type_entreprise", ""))
	setString(params, "date_radiation_rcs_minimum", getString(args, "date_radiation_rcs_minimum", ""))
	setString(params, "date_radiation_rcs_maximum", getString(args, "date_radiation_rcs_maximum", ""))
	setString(params, "capital_minimum", getString(args, "capital_minimum", ""))
	setString(params, "capital_maximum", getString(args, "capital_maximum", ""))
	setString(params, "chiffre_affaires_minimum", getString(args, "chiffre_affaires_minimum", ""))
	setString(params, "chiffre_affaires_maximum", getString(args, "chiffre_affaires_maximum", ""))
	setString(params, "resultat_minimum", getString(args, "resultat_minimum", ""))
	setString(params, "resultat_maximum", getString(args, "resultat_maximum", ""))
	return params
}

// extractFirstSiren parses a search response and extracts the SIREN of the first result.
func extractFirstSiren(data json.RawMessage) (string, error) {
	var result struct {
		Resultats []struct {
			Siren string `json:"siren"`
		} `json:"resultats"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("parsing search response: %w", err)
	}
	if len(result.Resultats) == 0 {
		return "", fmt.Errorf("no results found")
	}
	siren := result.Resultats[0].Siren
	if siren == "" {
		return "", fmt.Errorf("first result has no SIREN")
	}
	return siren, nil
}

// --- Schema helpers ---

// prop returns a JSON schema property definition as a map.
func prop(typ, description string) map[string]any {
	return map[string]any{
		"type":        typ,
		"description": description,
	}
}

// objectSchema returns a minimal JSON Schema object with the given properties and required fields.
func objectSchema(properties map[string]any, required []string) map[string]any {
	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

// commonSearchProperties returns shared schema properties for search endpoints.
func commonSearchProperties() map[string]any {
	return map[string]any{
		"q":                          prop("string", "Search query text"),
		"page":                       prop("integer", "Page number (starts at 1)"),
		"par_page":                   prop("integer", "Results per page (default 10, max 100)"),
		"precision":                  prop("string", "Search precision: standard, exacte, or phonetique"),
		"code_naf":                   prop("string", "NAF code filter"),
		"departement":                prop("string", "Department number filter"),
		"region":                     prop("string", "Region filter"),
		"code_postal":                prop("string", "Postal code filter"),
		"convention_collective":      prop("string", "Collective agreement filter"),
		"categorie_juridique":        prop("string", "Legal category code filter"),
		"entreprise_cessee":          prop("boolean", "Filter ceased companies (true=only ceased, false=only active)"),
		"statut_rcs":                 prop("string", "RCS status filter"),
		"objet_social":               prop("string", "Business object text filter"),
		"date_creation_minimum":      prop("string", "Minimum creation date (YYYY-MM-DD)"),
		"date_creation_maximum":      prop("string", "Maximum creation date (YYYY-MM-DD)"),
		"tranche_effectif":           prop("string", "Employee count bracket filter"),
		"type_entreprise":            prop("string", "Company type filter"),
		"date_radiation_rcs_minimum": prop("string", "Minimum RCS deregistration date (YYYY-MM-DD)"),
		"date_radiation_rcs_maximum": prop("string", "Maximum RCS deregistration date (YYYY-MM-DD)"),
		"capital_minimum":            prop("string", "Minimum share capital"),
		"capital_maximum":            prop("string", "Maximum share capital"),
		"chiffre_affaires_minimum":   prop("string", "Minimum revenue"),
		"chiffre_affaires_maximum":   prop("string", "Maximum revenue"),
		"resultat_minimum":           prop("string", "Minimum net income"),
		"resultat_maximum":           prop("string", "Maximum net income"),
	}
}

// directorSearchProperties returns additional schema properties for director search.
func directorSearchProperties() map[string]any {
	return map[string]any{
		"age_dirigeant_min":         prop("integer", "Minimum director age"),
		"age_dirigeant_max":         prop("integer", "Maximum director age"),
		"date_de_naissance_dirigeant_min": prop("string", "Minimum director birth date (YYYY-MM-DD)"),
		"date_de_naissance_dirigeant_max": prop("string", "Maximum director birth date (YYYY-MM-DD)"),
		"nationalite_dirigeant":     prop("string", "Director nationality filter"),
		"qualite_dirigeant":         prop("string", "Director role filter (e.g. Président, Directeur général)"),
	}
}

// beneficiarySearchProperties returns additional schema properties for beneficiary search.
func beneficiarySearchProperties() map[string]any {
	return map[string]any{
		"age_beneficiaire_min":      prop("integer", "Minimum beneficiary age"),
		"age_beneficiaire_max":      prop("integer", "Maximum beneficiary age"),
		"date_de_naissance_beneficiaire_min": prop("string", "Min beneficiary birth date (YYYY-MM-DD)"),
		"date_de_naissance_beneficiaire_max": prop("string", "Max beneficiary birth date (YYYY-MM-DD)"),
		"nationalite_beneficiaire":  prop("string", "Beneficiary nationality filter"),
		"type_beneficiaire":         prop("string", "Beneficiary type filter"),
	}
}

// mergeProperties merges multiple property maps into one.
func mergeProperties(maps ...map[string]any) map[string]any {
	result := make(map[string]any)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}
