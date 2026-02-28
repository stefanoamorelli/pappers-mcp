# Pappers MCP

<p align="center">
  <img alt="Go: 1.24+" src="https://img.shields.io/badge/go-1.24+-brightgreen.svg" />
  <a href="https://github.com/stefanoamorelli/pappers-mcp/actions/workflows/ci.yml"><img alt="CI" src="https://github.com/stefanoamorelli/pappers-mcp/actions/workflows/ci.yml/badge.svg" /></a>
  <img alt="License: AGPL-3.0" src="https://img.shields.io/badge/license-AGPL--3.0-blue.svg" />
  <br />
  <a href="https://modelcontextprotocol.io"><img alt="MCP" src="https://img.shields.io/badge/MCP-compatible-8A2BE2.svg" /></a>
  <a href="https://spec.modelcontextprotocol.io/specification/2025-03-26/"><img alt="MCP Spec" src="https://img.shields.io/badge/MCP_spec-2025--03--26-8A2BE2.svg" /></a>
  <img alt="Transport: stdio" src="https://img.shields.io/badge/transport-stdio-8A2BE2.svg" />
  <img alt="Tools: 24" src="https://img.shields.io/badge/tools-24-8A2BE2.svg" />
  <br />
  <a href="https://www.pappers.fr/api"><img alt="Pappers API v2" src="https://img.shields.io/badge/Pappers_API-v2-orange.svg" /></a>
</p>

MCP server for accessing the [Pappers API v2](https://www.pappers.fr/api). Connects AI assistants to French company data: legal info, financials, directors, beneficial owners, official documents, and more.

> [!IMPORTANT]
> This is an unofficial, community-built project. It is not affiliated with, endorsed by, or supported by [Pappers](https://www.pappers.fr). You must have a valid Pappers API key and comply with the [Pappers Terms of Use](https://www.pappers.fr/cgu) and [API documentation](https://www.pappers.fr/api/documentation) when using this server. You are responsible for your own API usage and any associated costs.

## Quick Start

### Claude Code

```bash
claude mcp add --transport stdio -e PAPPERS_API_KEY=your-key pappers-mcp \
  -- go run github.com/stefanoamorelli/pappers-mcp/cmd/pappers-mcp@latest
```

### Manual configuration

```json
{
  "mcpServers": {
    "pappers-mcp": {
      "command": "go",
      "args": ["run", "github.com/stefanoamorelli/pappers-mcp/cmd/pappers-mcp@latest"],
      "env": {
        "PAPPERS_API_KEY": "your-api-key"
      }
    }
  }
}
```

Or build from source:

```bash
go build -o pappers-mcp ./cmd/pappers-mcp
```

## Tools

| Category | Tools | Source |
|----------|-------|--------|
| **Company** | Company data (SIREN/SIRET/name), association data | [1], [2] |
| **Search** | Companies, directors, beneficial owners, documents, BODACC publications | [3] |
| **Suggestions** | Company name autocomplete | [4] |
| **Financials** | Annual accounts, corporate map | [5], [6] |
| **Compliance** | PEP/sanctions screening | [7] |
| **Documents** | Pappers extract, INPI extract, INSEE notice, bylaws, UBO declaration, financial/non-financial reports, document download | [8] |
| **Surveillance** | Company/director watch lists, notifications | [9] |
| **Account** | API credit usage | [10] |

24 tools covering the full [Pappers API v2][11].

[1]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/company.go
[2]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/association.go
[3]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/search.go
[4]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/suggestions.go
[5]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/accounts.go
[6]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/cartography.go
[7]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/compliance.go
[8]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/documents.go
[9]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/surveillance.go
[10]: https://github.com/stefanoamorelli/pappers-mcp/blob/6b003fc691d9bf8c9ada02a6a0213fa74523e029/internal/tools/credits.go
[11]: https://www.pappers.fr/api/documentation

## Development

```bash
make build          # Build binary
make test           # Run unit tests (no API key needed)
make conformance    # Verify all 24 tools are registered
```

Integration tests require a valid API key:

```bash
PAPPERS_API_KEY=xxx make test-integration
```

## License

[AGPL-3.0](LICENSE)

Copyright (c) 2025 [Stefano Amorelli](https://amorelli.tech) (stefano@amorelli.tech)
