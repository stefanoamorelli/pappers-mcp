package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stefanoamorelli/pappers-mcp/internal/client"
	"github.com/stefanoamorelli/pappers-mcp/internal/tools"
)

func main() {
	apiKey := os.Getenv("PAPPERS_API_KEY")
	if apiKey == "" {
		log.Fatal("PAPPERS_API_KEY environment variable is required")
	}

	c := client.New(apiKey)
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "pappers-mcp",
		Version: "2.0.0",
	}, nil)

	filter := tools.NewToolFilter(
		os.Getenv("PAPPERS_ENABLED_TOOLS"),
		os.Getenv("PAPPERS_DISABLED_TOOLS"),
	)
	tools.RegisterAll(server, c, filter)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
