package cmd

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestMeetingsListParseCustomRangeFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "from and to",
			args: []string{"meetings", "list", "--from", "2026-04-15", "--to", "2026-04-17", "--json"},
		},
		{
			name: "legacy start and end aliases",
			args: []string{"meetings", "list", "--start", "2026-04-15", "--end", "2026-04-17", "--json"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &CLI{}
			parser, err := kong.New(cli,
				kong.Name("granola"),
				kong.Description("A CLI for Granola meeting notes"),
				kong.UsageOnError(),
				kong.Vars{"version": "test"},
			)
			if err != nil {
				t.Fatalf("create parser: %v", err)
			}

			if _, err := parser.Parse(tt.args); err != nil {
				t.Fatalf("parse args: %v", err)
			}

			if cli.Meetings.List.From != "2026-04-15" {
				t.Fatalf("from = %q, want %q", cli.Meetings.List.From, "2026-04-15")
			}
			if cli.Meetings.List.To != "2026-04-17" {
				t.Fatalf("to = %q, want %q", cli.Meetings.List.To, "2026-04-17")
			}
			if !cli.Meetings.List.JSON {
				t.Fatal("json flag was not set")
			}
		})
	}
}

func TestMeetingsListToolArgs(t *testing.T) {
	tests := []struct {
		name    string
		cmd     MeetingsListCmd
		want    map[string]any
		wantErr string
	}{
		{
			name: "custom bounds imply custom range",
			cmd: MeetingsListCmd{
				From: "2026-04-15",
				To:   "2026-04-17",
			},
			want: map[string]any{
				"time_range":   "custom",
				"custom_start": "2026-04-15",
				"custom_end":   "2026-04-17",
			},
		},
		{
			name: "custom bounds conflict with preset range",
			cmd: MeetingsListCmd{
				Range: "last_week",
				From:  "2026-04-15",
			},
			wantErr: "--from/--to can't be used with --range=last_week",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cmd.toolArgs()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q", tt.wantErr)
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("error = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("tool args: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("args len = %d, want %d", len(got), len(tt.want))
			}
			for key, want := range tt.want {
				if got[key] != want {
					t.Fatalf("args[%q] = %v, want %v", key, got[key], want)
				}
			}
		})
	}
}
