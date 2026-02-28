// Package testutil provides test utilities for the Pappers MCP server.
package testutil

import (
	"context"
	"encoding/json"
	"net/url"
)

// MockPappersClient is a mock implementation of client.PappersClient for testing.
type MockPappersClient struct {
	GetCompanyFunc              func(ctx context.Context, params url.Values) (json.RawMessage, error)
	GetAssociationFunc          func(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchCompaniesFunc         func(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchDirectorsFunc         func(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchBeneficiariesFunc     func(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchDocumentsFunc         func(ctx context.Context, params url.Values) (json.RawMessage, error)
	SearchPublicationsFunc      func(ctx context.Context, params url.Values) (json.RawMessage, error)
	SuggestFunc                 func(ctx context.Context, params url.Values) (json.RawMessage, error)
	GetAnnualAccountsFunc       func(ctx context.Context, params url.Values) (json.RawMessage, error)
	GetCorporateMapFunc         func(ctx context.Context, params url.Values) (json.RawMessage, error)
	CheckPEPFunc                func(ctx context.Context, params url.Values) (json.RawMessage, error)
	DownloadDocumentFunc        func(ctx context.Context, token string) ([]byte, string, error)
	DownloadCompanyDocumentFunc func(ctx context.Context, path string, params url.Values) ([]byte, string, error)
	AddCompanyWatchFunc         func(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error)
	AddDirectorWatchFunc        func(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error)
	DeleteNotificationsFunc     func(ctx context.Context, params url.Values) (json.RawMessage, error)
	AddNotificationInfoFunc     func(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error)
	GetAPICreditsFunc           func(ctx context.Context) (json.RawMessage, error)
}

func (m *MockPappersClient) GetCompany(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.GetCompanyFunc != nil {
		return m.GetCompanyFunc(ctx, params)
	}
	return json.RawMessage(`{}`), nil
}

func (m *MockPappersClient) GetAssociation(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.GetAssociationFunc != nil {
		return m.GetAssociationFunc(ctx, params)
	}
	return json.RawMessage(`{}`), nil
}

func (m *MockPappersClient) SearchCompanies(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.SearchCompaniesFunc != nil {
		return m.SearchCompaniesFunc(ctx, params)
	}
	return json.RawMessage(`{"resultats":[]}`), nil
}

func (m *MockPappersClient) SearchDirectors(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.SearchDirectorsFunc != nil {
		return m.SearchDirectorsFunc(ctx, params)
	}
	return json.RawMessage(`{"resultats":[]}`), nil
}

func (m *MockPappersClient) SearchBeneficiaries(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.SearchBeneficiariesFunc != nil {
		return m.SearchBeneficiariesFunc(ctx, params)
	}
	return json.RawMessage(`{"resultats":[]}`), nil
}

func (m *MockPappersClient) SearchDocuments(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.SearchDocumentsFunc != nil {
		return m.SearchDocumentsFunc(ctx, params)
	}
	return json.RawMessage(`{"resultats":[]}`), nil
}

func (m *MockPappersClient) SearchPublications(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.SearchPublicationsFunc != nil {
		return m.SearchPublicationsFunc(ctx, params)
	}
	return json.RawMessage(`{"resultats":[]}`), nil
}

func (m *MockPappersClient) Suggest(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.SuggestFunc != nil {
		return m.SuggestFunc(ctx, params)
	}
	return json.RawMessage(`[]`), nil
}

func (m *MockPappersClient) GetAnnualAccounts(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.GetAnnualAccountsFunc != nil {
		return m.GetAnnualAccountsFunc(ctx, params)
	}
	return json.RawMessage(`{}`), nil
}

func (m *MockPappersClient) GetCorporateMap(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.GetCorporateMapFunc != nil {
		return m.GetCorporateMapFunc(ctx, params)
	}
	return json.RawMessage(`{}`), nil
}

func (m *MockPappersClient) CheckPEP(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.CheckPEPFunc != nil {
		return m.CheckPEPFunc(ctx, params)
	}
	return json.RawMessage(`{}`), nil
}

func (m *MockPappersClient) DownloadDocument(ctx context.Context, token string) ([]byte, string, error) {
	if m.DownloadDocumentFunc != nil {
		return m.DownloadDocumentFunc(ctx, token)
	}
	return []byte("PDF"), "application/pdf", nil
}

func (m *MockPappersClient) DownloadCompanyDocument(ctx context.Context, path string, params url.Values) ([]byte, string, error) {
	if m.DownloadCompanyDocumentFunc != nil {
		return m.DownloadCompanyDocumentFunc(ctx, path, params)
	}
	return []byte("PDF"), "application/pdf", nil
}

func (m *MockPappersClient) AddCompanyWatch(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
	if m.AddCompanyWatchFunc != nil {
		return m.AddCompanyWatchFunc(ctx, listID, body)
	}
	return json.RawMessage(`{"success":true}`), nil
}

func (m *MockPappersClient) AddDirectorWatch(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
	if m.AddDirectorWatchFunc != nil {
		return m.AddDirectorWatchFunc(ctx, listID, body)
	}
	return json.RawMessage(`{"success":true}`), nil
}

func (m *MockPappersClient) DeleteNotifications(ctx context.Context, params url.Values) (json.RawMessage, error) {
	if m.DeleteNotificationsFunc != nil {
		return m.DeleteNotificationsFunc(ctx, params)
	}
	return json.RawMessage(`{"success":true}`), nil
}

func (m *MockPappersClient) AddNotificationInfo(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
	if m.AddNotificationInfoFunc != nil {
		return m.AddNotificationInfoFunc(ctx, listID, body)
	}
	return json.RawMessage(`{"success":true}`), nil
}

func (m *MockPappersClient) GetAPICredits(ctx context.Context) (json.RawMessage, error) {
	if m.GetAPICreditsFunc != nil {
		return m.GetAPICreditsFunc(ctx)
	}
	return json.RawMessage(`{"jetons_utilises":100,"jetons_restants":900}`), nil
}
