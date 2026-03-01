package cmd

import (
	"context"
	"fmt"

	"github.com/lox/granola-cli/internal/cli"
	"github.com/lox/granola-cli/internal/output"
)

type MeetingsCmd struct {
	List       MeetingsListCmd       `cmd:"" help:"List meetings"`
	View       MeetingsViewCmd       `cmd:"" help:"View meeting details"`
	Transcript MeetingsTranscriptCmd `cmd:"" help:"View meeting transcript"`
}

type MeetingsListCmd struct {
	Range string `help:"Time range: this_week, last_week, last_30_days, or custom" short:"r" default:"this_week" enum:"this_week,last_week,last_30_days,custom"`
	Start string `help:"Custom range start (ISO date, requires --range=custom)" short:"s"`
	End   string `help:"Custom range end (ISO date, requires --range=custom)" short:"e"`
	JSON  bool   `help:"Output as JSON" short:"j"`
}

func (c *MeetingsListCmd) Run(ctx *Context) error {
	ctx.JSON = c.JSON

	client, err := cli.RequireClient()
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	bgCtx := context.Background()

	args := map[string]any{
		"time_range": c.Range,
	}
	if c.Range == "custom" {
		if c.Start != "" {
			args["custom_start"] = c.Start
		}
		if c.End != "" {
			args["custom_end"] = c.End
		}
	}

	result, err := client.CallToolText(bgCtx, "list_meetings", args)
	if err != nil {
		output.PrintError(err)
		return err
	}

	meetings, err := output.ParseMeetingsList(result)
	if err != nil {
		fmt.Println(result)
		return nil
	}

	return output.PrintMeetings(meetings, ctx.JSON)
}

type MeetingsViewCmd struct {
	ID   string `arg:"" help:"Meeting ID (UUID)"`
	JSON bool   `help:"Output as JSON" short:"j"`
	Raw  bool   `help:"Output raw response without formatting" short:"r"`
}

func (c *MeetingsViewCmd) Run(ctx *Context) error {
	ctx.JSON = c.JSON

	client, err := cli.RequireClient()
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	bgCtx := context.Background()

	result, err := client.CallToolText(bgCtx, "get_meetings", map[string]any{
		"meeting_ids": []string{c.ID},
	})
	if err != nil {
		output.PrintError(err)
		return err
	}

	if c.Raw {
		fmt.Println(result)
		return nil
	}

	meeting, err := output.ParseMeetingDetail(result)
	if err != nil {
		fmt.Println(result)
		return nil
	}

	return output.PrintMeetingDetail(meeting, ctx.JSON)
}

type MeetingsTranscriptCmd struct {
	ID  string `arg:"" help:"Meeting ID (UUID)"`
	Raw bool   `help:"Output raw response without formatting" short:"r"`
}

func (c *MeetingsTranscriptCmd) Run(ctx *Context) error {
	client, err := cli.RequireClient()
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()

	bgCtx := context.Background()

	result, err := client.CallToolText(bgCtx, "get_meeting_transcript", map[string]any{
		"meeting_id": c.ID,
	})
	if err != nil {
		output.PrintError(err)
		return err
	}

	if c.Raw {
		fmt.Println(result)
		return nil
	}

	return output.PrintTranscript(result)
}
