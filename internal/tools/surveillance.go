package tools

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
)

// --- Add Company Watch ---

func addCompanyWatchTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "add_company_watch",
		Description: "Add a company to a surveillance/monitoring list. You will be notified of changes to the company (e.g. director changes, financial updates, legal events).",
		InputSchema: objectSchema(map[string]any{
			"id_liste": prop("string", "ID of the surveillance list to add the company to"),
			"siren":    prop("string", "SIREN number of the company to watch (9 digits)"),
		}, []string{"siren"}),
	}
}

func addCompanyWatchHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		siren := getString(args, "siren", "")
		if siren == "" {
			return toolError("Parameter 'siren' is required"), nil
		}

		listID := getString(args, "id_liste", "")
		body, _ := json.Marshal(map[string]string{"siren": siren})

		data, err := c.AddCompanyWatch(ctx, listID, body)
		if err != nil {
			return toolErrorf("Failed to add company watch: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}

// --- Add Director Watch ---

func addDirectorWatchTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "add_director_watch",
		Description: "Add a director/person to a surveillance/monitoring list. You will be notified when this person's directorships or company roles change.",
		InputSchema: objectSchema(map[string]any{
			"id_liste":        prop("string", "ID of the surveillance list"),
			"nom":             prop("string", "Last name of the director"),
			"prenom":          prop("string", "First name of the director"),
			"date_de_naissance": prop("string", "Birth date (YYYY-MM-DD)"),
		}, []string{"nom", "prenom"}),
	}
}

func addDirectorWatchHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		nom := getString(args, "nom", "")
		prenom := getString(args, "prenom", "")
		if nom == "" || prenom == "" {
			return toolError("Parameters 'nom' and 'prenom' are required"), nil
		}

		listID := getString(args, "id_liste", "")
		bodyMap := map[string]string{
			"nom":    nom,
			"prenom": prenom,
		}
		if dob := getString(args, "date_de_naissance", ""); dob != "" {
			bodyMap["date_de_naissance"] = dob
		}
		body, _ := json.Marshal(bodyMap)

		data, err := c.AddDirectorWatch(ctx, listID, body)
		if err != nil {
			return toolErrorf("Failed to add director watch: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}

// --- Delete Notifications ---

func deleteNotificationsTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "delete_notifications",
		Description: "Delete surveillance notifications. Remove notifications you've already processed or that are no longer relevant.",
		InputSchema: objectSchema(map[string]any{
			"id_liste":       prop("string", "ID of the surveillance list"),
			"id_notification": prop("string", "ID of a specific notification to delete"),
			"all":            prop("boolean", "Set to true to delete all notifications"),
		}, nil),
	}
}

func deleteNotificationsHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		params := url.Values{}
		setString(params, "id_liste", getString(args, "id_liste", ""))
		setString(params, "id_notification", getString(args, "id_notification", ""))
		setBool(params, "all", args)

		data, err := c.DeleteNotifications(ctx, params)
		if err != nil {
			return toolErrorf("Failed to delete notifications: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}

// --- Add Notification Info ---

func addNotificationInfoTool() *mcp.Tool {
	return &mcp.Tool{
		Name:        "add_notification_info",
		Description: "Add or update information for a surveillance notification list (e.g. webhook URL, email address for alerts).",
		InputSchema: objectSchema(map[string]any{
			"id_liste": prop("string", "ID of the surveillance list"),
			"email":    prop("string", "Email address for notifications"),
			"webhook":  prop("string", "Webhook URL for notifications"),
		}, nil),
	}
}

func addNotificationInfoHandler(c client.PappersClient) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := extractArgs(req)
		if err != nil {
			return toolError(err.Error()), nil
		}

		listID := getString(args, "id_liste", "")
		bodyMap := make(map[string]string)
		if email := getString(args, "email", ""); email != "" {
			bodyMap["email"] = email
		}
		if webhook := getString(args, "webhook", ""); webhook != "" {
			bodyMap["webhook"] = webhook
		}
		body, _ := json.Marshal(bodyMap)

		data, err := c.AddNotificationInfo(ctx, listID, body)
		if err != nil {
			return toolErrorf("Failed to update notification info: %v", err), nil
		}

		return toolText(string(data)), nil
	}
}
