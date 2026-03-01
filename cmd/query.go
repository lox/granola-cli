package cmd

import (
	"context"
	"fmt"

	"github.com/lox/granola-cli/internal/cli"
	"github.com/lox/granola-cli/internal/output"
)

type QueryCmd struct {
	Query      string   `arg:"" help:"Natural language query about your meetings"`
	MeetingIDs []string `help:"Limit context to specific meeting IDs" short:"m" name:"meeting-id"`
	Raw        bool     `help:"Output raw response without formatting" short:"r"`
}

func (c *QueryCmd) Run(ctx *Context) error {
	client, err := cli.RequireClient()
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	bgCtx := context.Background()

	args := map[string]any{
		"query": c.Query,
	}
	if len(c.MeetingIDs) > 0 {
		args["document_ids"] = c.MeetingIDs
	}

	result, err := client.CallToolText(bgCtx, "query_granola_meetings", args)
	if err != nil {
		output.PrintError(err)
		return err
	}

	if c.Raw {
		fmt.Println(result)
		return nil
	}

	return output.RenderMarkdown(result)
}
