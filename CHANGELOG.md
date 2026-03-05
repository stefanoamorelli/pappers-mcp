# Changelog

## Unreleased

## v0.2.0 (2026-03-05)

- Switch from stdio to streamable HTTP transport
- Add `/health` endpoint with status and version info
- Add `/mcp` endpoint for MCP protocol over HTTP
- Server port configurable via `PORT` environment variable (default `3000`)
- Add tool filtering via `PAPPERS_ENABLED_TOOLS` and `PAPPERS_DISABLED_TOOLS` environment variables
- Add release workflow for cross-compiled binaries on version tags
- Add CHANGELOG.md

## v0.1.0 (2026-02-28)

- Initial release with 24 Pappers API tools
- MCP server with stdio transport
- CI workflow
