---
name: granola
description: Search and query Granola meeting notes from the command line. List meetings, view summaries, read transcripts, and ask natural language questions about meeting content.
allowed-tools: Bash(granola-cli:*)
---

# Granola CLI

A CLI to access Granola meeting notes from the command line, using Granola's remote MCP server.

## Prerequisites

The `granola-cli` command must be available on PATH. To check:

```bash
granola-cli --version
```

If not installed:

```bash
go install github.com/lox/granola-cli@latest
```

Or see: https://github.com/lox/granola-cli

## Authentication

The CLI uses OAuth authentication. On first use, it opens a browser for authorization:

```bash
granola-cli auth login      # Authenticate with Granola
granola-cli auth status     # Check authentication status
granola-cli auth refresh    # Refresh access token
granola-cli auth logout     # Clear credentials
```

For CI/headless environments, set `GRANOLA_ACCESS_TOKEN` environment variable.

## Available Commands

```
granola-cli auth            # Manage authentication
granola-cli meetings        # List, view, and get transcripts
granola-cli query           # Natural language query about meetings
granola-cli tools           # List available MCP tools
```

## Common Operations

### List Meetings

```bash
granola-cli meetings list                        # This week (default)
granola-cli meetings list -r last_week           # Last week
granola-cli meetings list -r last_30_days        # Last 30 days
granola-cli meetings list -r custom -s 2026-01-01 -e 2026-01-31  # Custom range
granola-cli meetings list --json                 # JSON output
```

### View Meeting Details

Shows title, date, attendees, AI summary, and notes.

```bash
granola-cli meetings view <meeting-id>           # View meeting details
granola-cli meetings view <meeting-id> --raw     # Raw MCP response
granola-cli meetings view <meeting-id> --json    # JSON output
```

### Read Transcripts

```bash
granola-cli meetings transcript <meeting-id>     # View transcript
granola-cli meetings transcript <meeting-id> --raw  # Raw response
```

### Natural Language Query

Ask questions about your meetings. Returns synthesised answers with citation links.

```bash
granola-cli query "what was discussed about hiring last week"
granola-cli query "action items from my 1:1s"
granola-cli query "what decisions were made about the roadmap"
granola-cli query "summarise my meetings with Kevin" --raw
granola-cli query "test engine updates" -m <meeting-id>   # Limit to specific meetings
```

## Output Formats

Most commands support `--json` for machine-readable output:

```bash
granola-cli meetings list --json | jq '.[0].id'
granola-cli meetings list -r last_week --json | jq '.[].title'
```

## Tips for Agents

1. **Use `query` for content questions** — it searches across all meetings and returns synthesised answers with citations
2. **Use `meetings list` then `meetings view`** — to browse meetings and drill into specifics
3. **Use transcripts for exact quotes** — `meetings transcript` returns verbatim text, `meetings view` returns AI summaries
4. **Copy full UUIDs from list output** — the table shows truncated 8-char IDs, use `--json` to get full UUIDs
5. **Check --help** — every command has detailed help: `granola-cli meetings list --help`
6. **Raw output for debugging** — use `--raw` to see the original MCP response
