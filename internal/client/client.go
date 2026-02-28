// Package client provides an HTTP client for the Pappers API v2.
package client

import (
	stdbytes "bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://api.pappers.fr/v2"

// PappersClient defines the interface for interacting with the Pappers API.
// JSON endpoints return json.RawMessage so that MCP tools can pass the raw
// response through to the LLM without mapping into Go structs.
type PappersClient interface {
	// Company data
	GetCompany(ctx context.Context, params url.Values) (json.RawMessage, error)
	GetAssociation(ctx context.Context, params url.Values) (json.RawMessage, error)

	// Search
	SearchCompanies(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchDirectors(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchBeneficiaries(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchDocuments(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchPublications(ctx context.Context, params url.Values) (json.RawMessage, error)

	// Suggestions
	Suggest(ctx context.Context, params url.Values) (json.RawMessage, error)

	// Accounts
	GetAnnualAccounts(ctx context.Context, params url.Values) (json.RawMessage, error)

	// Cartography
	GetCorporateMap(ctx context.Context, params url.Values) (json.RawMessage, error)

	// Compliance
	CheckPEP(ctx context.Context, params url.Values) (json.RawMessage, error)

	// Document downloads (binary)
	DownloadDocument(ctx context.Context, token string) ([]byte, string, error)
	DownloadCompanyDocument(ctx context.Context, path string, params url.Values) ([]byte, string, error)

	// Surveillance
	AddCompanyWatch(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error)
	AddDirectorWatch(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error)
	DeleteNotifications(ctx context.Context, params url.Values) (json.RawMessage, error)
	AddNotificationInfo(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error)

	// Credits
	GetAPICredits(ctx context.Context) (json.RawMessage, error)
}

// Option configures the HTTP client.
type Option func(*httpPappersClient)

// WithBaseURL overrides the default Pappers API base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *httpPappersClient) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *httpPappersClient) {
		c.httpClient = hc
	}
}

// New creates a new PappersClient with the given API token and options.
func New(apiToken string, opts ...Option) PappersClient {
	c := &httpPappersClient{
		apiToken:   apiToken,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type httpPappersClient struct {
	apiToken   string
	baseURL    string
	httpClient *http.Client
}

// doJSON performs a GET request and returns the raw JSON response body.
func (c *httpPappersClient) doJSON(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_token", c.apiToken)

	reqURL := c.baseURL + path + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return json.RawMessage(body), nil
}

// doBinary performs a GET request and returns the raw bytes and content type.
func (c *httpPappersClient) doBinary(ctx context.Context, path string, params url.Values) ([]byte, string, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_token", c.apiToken)

	reqURL := c.baseURL + path + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, resp.Header.Get("Content-Type"), nil
}

// doPost performs a POST request with JSON body and returns the raw JSON response.
func (c *httpPappersClient) doPost(ctx context.Context, path string, params url.Values, body io.Reader) (json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_token", c.apiToken)

	reqURL := c.baseURL + path + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

// doDelete performs a DELETE request and returns the raw JSON response.
func (c *httpPappersClient) doDelete(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("api_token", c.apiToken)

	reqURL := c.baseURL + path + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return json.RawMessage(body), nil
}

// --- Interface method implementations ---

func (c *httpPappersClient) GetCompany(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/entreprise", params)
}

func (c *httpPappersClient) GetAssociation(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/association", params)
}

func (c *httpPappersClient) SearchCompanies(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/recherche", params)
}

func (c *httpPappersClient) SearchDirectors(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/recherche-dirigeants", params)
}

func (c *httpPappersClient) SearchBeneficiaries(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/recherche-beneficiaires", params)
}

func (c *httpPappersClient) SearchDocuments(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/recherche-documents", params)
}

func (c *httpPappersClient) SearchPublications(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/recherche-publications", params)
}

func (c *httpPappersClient) Suggest(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/suggestions", params)
}

func (c *httpPappersClient) GetAnnualAccounts(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/entreprise/comptes", params)
}

func (c *httpPappersClient) GetCorporateMap(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/cartographie", params)
}

func (c *httpPappersClient) CheckPEP(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doJSON(ctx, "/conformite/personne_politiquement_exposee", params)
}

func (c *httpPappersClient) DownloadDocument(ctx context.Context, token string) ([]byte, string, error) {
	params := url.Values{}
	params.Set("token", token)
	return c.doBinary(ctx, "/document/telechargement", params)
}

func (c *httpPappersClient) DownloadCompanyDocument(ctx context.Context, path string, params url.Values) ([]byte, string, error) {
	return c.doBinary(ctx, path, params)
}

func (c *httpPappersClient) AddCompanyWatch(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
	params := url.Values{}
	if listID != "" {
		params.Set("id_liste", listID)
	}
	return c.doPost(ctx, "/surveillance/entreprise", params, bodyReader(body))
}

func (c *httpPappersClient) AddDirectorWatch(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
	params := url.Values{}
	if listID != "" {
		params.Set("id_liste", listID)
	}
	return c.doPost(ctx, "/surveillance/dirigeant", params, bodyReader(body))
}

func (c *httpPappersClient) DeleteNotifications(ctx context.Context, params url.Values) (json.RawMessage, error) {
	return c.doDelete(ctx, "/surveillance/notifications", params)
}

func (c *httpPappersClient) AddNotificationInfo(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
	params := url.Values{}
	if listID != "" {
		params.Set("id_liste", listID)
	}
	return c.doPost(ctx, "/surveillance/liste/informations", params, bodyReader(body))
}

func (c *httpPappersClient) GetAPICredits(ctx context.Context) (json.RawMessage, error) {
	return c.doJSON(ctx, "/suivi-jetons", nil)
}

// bodyReader wraps json.RawMessage in a reader, handling nil gracefully.
func bodyReader(data json.RawMessage) io.Reader {
	if data == nil {
		return nil
	}
	return stdbytes.NewReader(data)
}
