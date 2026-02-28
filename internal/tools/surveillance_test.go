package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stefanoamorelli/pappers-mcp/internal/testutil"
)

func TestAddCompanyWatchHandler_RequiresSiren(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := addCompanyWatchHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when siren is missing")
	}
}

func TestAddCompanyWatchHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		AddCompanyWatchFunc: func(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
			if listID != "list1" {
				t.Errorf("expected listID=list1, got %s", listID)
			}
			return json.RawMessage(`{"success":true}`), nil
		},
	}

	handler := addCompanyWatchHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"siren":    "443061841",
		"id_liste": "list1",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestAddDirectorWatchHandler_RequiresNomPrenom(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := addDirectorWatchHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"nom": "Doe"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Fatal("expected error when prenom is missing")
	}
}

func TestAddDirectorWatchHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{
		AddDirectorWatchFunc: func(ctx context.Context, listID string, body json.RawMessage) (json.RawMessage, error) {
			return json.RawMessage(`{"success":true}`), nil
		},
	}

	handler := addDirectorWatchHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"nom":    "Doe",
		"prenom": "John",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestDeleteNotificationsHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := deleteNotificationsHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{"all": true}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}

func TestAddNotificationInfoHandler_HappyPath(t *testing.T) {
	mock := &testutil.MockPappersClient{}
	handler := addNotificationInfoHandler(mock)
	result, err := handler(context.Background(), makeReq(map[string]any{
		"id_liste": "list1",
		"email":    "test@example.com",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal("expected success")
	}
}
