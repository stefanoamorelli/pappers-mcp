package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
	"github.com/stefanoamorelli/pappers-mcp/internal/tools"
)

const version = "0.2.0"

func main() {
	apiKey := os.Getenv("PAPPERS_API_KEY")
	if apiKey == "" {
		log.Fatal("PAPPERS_API_KEY environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	c := client.New(apiKey)
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "pappers-mcp",
		Version: version,
	}, nil)

	filter := tools.NewToolFilter(
		os.Getenv("PAPPERS_ENABLED_TOOLS"),
		os.Getenv("PAPPERS_DISABLED_TOOLS"),
	)
	tools.RegisterAll(server, c, filter)

	mcpHandler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"version": version,
		})
	})

	addr := net.JoinHostPort("", port)
	log.Printf("pappers-mcp %s listening on %s", version, addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
