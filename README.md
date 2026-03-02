# granola-cli

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.25-blue)](https://golang.org/)

A command-line interface for [Granola](https://www.granola.ai) meeting notes using the remote MCP (Model Context Protocol).

**Works great with AI agents** — includes a [skill](#skills) that lets agents search and query your meeting notes alongside your code.

## Installation

### From Source

```bash
go install github.com/lox/granola-cli@latest
```

### Build Locally

```bash
git clone https://github.com/lox/granola-cli
cd granola-cli
mise run build
```

## Quick Start

```bash
# Authenticate with Granola (opens browser for OAuth)
granola-cli auth login

# List this week's meetings
granola-cli meetings list

# View a meeting's summary and notes
granola-cli meetings view <meeting-id>

# Ask a question about your meetings
granola-cli query "what was discussed about hiring last week"
```

## Commands

### Authentication

```bash
granola-cli auth login      # Authenticate with Granola via OAuth
granola-cli auth refresh    # Refresh the access token
granola-cli auth status     # Show authentication status
granola-cli auth logout     # Clear stored credentials
```

### Meetings

```bash
granola-cli meetings list                        # List this week's meetings (default)
granola-cli meetings list -r last_week           # Last week
granola-cli meetings list -r last_30_days        # Last 30 days
granola-cli meetings list -r custom -s 2026-01-01 -e 2026-01-31  # Custom range
granola-cli meetings list --json                 # Output as JSON

granola-cli meetings view <meeting-id>           # View meeting details
granola-cli meetings view <meeting-id> --raw     # Raw MCP response
granola-cli meetings view <meeting-id> --json    # Output as JSON

granola-cli meetings transcript <meeting-id>     # View full transcript
granola-cli meetings transcript <meeting-id> --raw
```

### Query

Ask natural language questions about your meetings. Returns synthesised answers with citation links.

```bash
granola-cli query "what was discussed about hiring last week"
granola-cli query "action items from my 1:1s"
granola-cli query "what decisions were made about the roadmap"
granola-cli query "summarise my meetings with Kevin" --raw
granola-cli query "test engine updates" -m <meeting-id>   # Limit to specific meetings
```

### Other

```bash
granola-cli tools                              # List available MCP tools
granola-cli version                            # Show version
granola-cli --help                             # Show help
```

## Configuration

Credentials are stored at `~/.config/granola-cli/token.json`.

The CLI uses Granola's remote MCP server with OAuth authentication. On first run, `granola-cli auth login` will open your browser to authorise the CLI with your Granola account.

**Note:** Access tokens expire periodically. The CLI automatically refreshes tokens when they expire or are about to expire, so you typically don't need to think about this. Use `granola-cli auth refresh` to manually refresh if needed.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GRANOLA_ACCESS_TOKEN` | Access token for CI/headless usage (skips OAuth) |

## How It Works

This CLI connects to [Granola's remote MCP server](https://www.granola.ai/blog/granola-mcp) at `https://mcp.granola.ai/mcp` using the Model Context Protocol. This provides:

- **OAuth authentication** — No API tokens to manage
- **Meeting summaries** — AI-generated summaries and notes
- **Full transcripts** — Verbatim meeting transcripts
- **Natural language queries** — Ask questions across all your meetings

## Skills

granola-cli includes a skill that helps AI agents use the CLI effectively.

### Amp / Claude Code

Install the skill using [skills.sh](https://skills.sh):

```bash
npx skills add lox/granola-cli
```

View the skill at: [skills/granola/SKILL.md](skills/granola/SKILL.md)

## Links

- [Granola MCP Announcement](https://www.granola.ai/blog/granola-mcp)
- [Model Context Protocol](https://modelcontextprotocol.io/)
