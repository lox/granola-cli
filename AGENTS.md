# Agent Instructions

## Project

granola-cli is a Go CLI that wraps the Granola MCP server at `https://mcp.granola.ai/mcp`. It follows the same patterns as [notion-cli](https://github.com/lox/notion-cli).

## Structure

- `cmd/` — Kong CLI commands (auth, meetings, query, tools)
- `internal/mcp/` — MCP client, OAuth flow, token storage
- `internal/cli/` — Client context and auto-refresh
- `internal/output/` — Formatting, tables, markdown rendering, XML parsing
- `skills/` — Amp skill definition

## Conventions

- Uses [Kong](https://github.com/alecthomas/kong) for CLI parsing
- Uses [mcp-go](https://github.com/mark3labs/mcp-go) for MCP client
- Uses [glamour](https://github.com/charmbracelet/glamour) for terminal markdown rendering
- OAuth tokens stored in `~/.config/granola-cli/token.json`
- Granola MCP returns XML responses that need sanitisation (unescaped `<>`, `&`)
- Conventional commits (`feat:`, `fix:`, etc.)

## Build & Test

```bash
mise run build    # Build binary
mise run test     # Run tests
mise run lint     # Run linter
mise run check    # All checks
```
