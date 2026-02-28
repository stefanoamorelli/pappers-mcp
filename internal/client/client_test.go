package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDoJSON_InjectsAPIToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("api_token")
		if token != "test-key" {
			t.Errorf("expected api_token=test-key, got %q", token)
		}
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New("test-key", WithBaseURL(srv.URL))
	data, err := c.GetCompany(context.Background(), url.Values{"siren": {"123456789"}})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"ok":true}` {
		t.Errorf("unexpected response: %s", data)
	}
}

func TestDoJSON_CorrectPath(t *testing.T) {
	tests := []struct {
		name   string
		call   func(PappersClient, context.Context) (json.RawMessage, error)
		path   string
	}{
		{"GetCompany", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.GetCompany(ctx, url.Values{"siren": {"123"}})
		}, "/entreprise"},
		{"GetAssociation", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.GetAssociation(ctx, url.Values{"siren": {"123"}})
		}, "/association"},
		{"SearchCompanies", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.SearchCompanies(ctx, url.Values{"q": {"test"}})
		}, "/recherche"},
		{"SearchDirectors", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.SearchDirectors(ctx, url.Values{"q": {"test"}})
		}, "/recherche-dirigeants"},
		{"SearchBeneficiaries", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.SearchBeneficiaries(ctx, url.Values{"q": {"test"}})
		}, "/recherche-beneficiaires"},
		{"SearchDocuments", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.SearchDocuments(ctx, url.Values{"q": {"test"}})
		}, "/recherche-documents"},
		{"SearchPublications", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.SearchPublications(ctx, url.Values{"q": {"test"}})
		}, "/recherche-publications"},
		{"Suggest", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.Suggest(ctx, url.Values{"q": {"goo"}})
		}, "/suggestions"},
		{"GetAnnualAccounts", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.GetAnnualAccounts(ctx, url.Values{"siren": {"123"}})
		}, "/entreprise/comptes"},
		{"GetCorporateMap", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.GetCorporateMap(ctx, url.Values{"siren": {"123"}})
		}, "/cartographie"},
		{"CheckPEP", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.CheckPEP(ctx, url.Values{"nom": {"Doe"}, "prenom": {"John"}})
		}, "/conformite/personne_politiquement_exposee"},
		{"GetAPICredits", func(c PappersClient, ctx context.Context) (json.RawMessage, error) {
			return c.GetAPICredits(ctx)
		}, "/suivi-jetons"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Errorf("expected path %s, got %s", tt.path, r.URL.Path)
				}
				w.Write([]byte(`{}`))
			}))
			defer srv.Close()

			c := New("key", WithBaseURL(srv.URL))
			_, err := tt.call(c, context.Background())
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestDoJSON_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	c := New("key", WithBaseURL(srv.URL))
	_, err := c.GetCompany(context.Background(), url.Values{"siren": {"000"}})
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestDoJSON_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`internal error`))
	}))
	defer srv.Close()

	c := New("key", WithBaseURL(srv.URL))
	_, err := c.SearchCompanies(context.Background(), url.Values{"q": {"test"}})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestDownloadDocument_CorrectTokenParam(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/document/telechargement" {
			t.Errorf("expected path /document/telechargement, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("token") != "my-token" {
			t.Errorf("expected token=my-token, got %s", r.URL.Query().Get("token"))
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("%PDF-1.4"))
	}))
	defer srv.Close()

	c := New("key", WithBaseURL(srv.URL))
	data, ct, err := c.DownloadDocument(context.Background(), "my-token")
	if err != nil {
		t.Fatal(err)
	}
	if ct != "application/pdf" {
		t.Errorf("expected application/pdf, got %s", ct)
	}
	if string(data) != "%PDF-1.4" {
		t.Errorf("unexpected data: %s", data)
	}
}

func TestDownloadCompanyDocument_CorrectPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/document/extrait_pappers" {
			t.Errorf("expected path /document/extrait_pappers, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("siren") != "443061841" {
			t.Errorf("expected siren param, got %s", r.URL.Query().Get("siren"))
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Write([]byte("data"))
	}))
	defer srv.Close()

	c := New("key", WithBaseURL(srv.URL))
	_, ct, err := c.DownloadCompanyDocument(context.Background(), "/document/extrait_pappers", url.Values{"siren": {"443061841"}})
	if err != nil {
		t.Fatal(err)
	}
	if ct != "application/pdf" {
		t.Errorf("expected application/pdf, got %s", ct)
	}
}

func TestAddCompanyWatch_PostMethod(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/surveillance/entreprise" {
			t.Errorf("expected path /surveillance/entreprise, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"success":true}`))
	}))
	defer srv.Close()

	c := New("key", WithBaseURL(srv.URL))
	data, err := c.AddCompanyWatch(context.Background(), "list1", json.RawMessage(`{"siren":"123"}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"success":true}` {
		t.Errorf("unexpected response: %s", data)
	}
}

func TestDeleteNotifications_DeleteMethod(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/surveillance/notifications" {
			t.Errorf("expected path /surveillance/notifications, got %s", r.URL.Path)
		}
		w.Write([]byte(`{"deleted":true}`))
	}))
	defer srv.Close()

	c := New("key", WithBaseURL(srv.URL))
	_, err := c.DeleteNotifications(context.Background(), url.Values{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 5000000000} // 5 seconds
	c := New("key", WithHTTPClient(customClient))
	hc := c.(*httpPappersClient)
	if hc.httpClient != customClient {
		t.Error("expected custom HTTP client to be set")
	}
}
