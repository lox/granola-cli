package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

func PrintMeetings(meetings []Meeting, asJSON bool) error {
	if asJSON {
		return printJSON(meetings)
	}

	if len(meetings) == 0 {
		fmt.Println("No meetings found.")
		return nil
	}

	table := NewTable("ID", "TITLE", "DATE", "TIME")
	for _, m := range meetings {
		table.AddRow(
			TruncateID(m.ID),
			Truncate(m.Title, 50),
			formatDate(m.StartTime),
			m.StartTime.Format("15:04"),
		)
	}
	table.Render()
	return nil
}

func PrintMeetingDetail(meeting *MeetingDetail, asJSON bool) error {
	if asJSON {
		return printJSON(meeting)
	}

	titleStyle := color.New(color.Bold, color.FgWhite)
	labelStyle := color.New(color.Faint)

	_, _ = titleStyle.Println(meeting.Title)
	fmt.Println()

	_, _ = labelStyle.Print("ID:        ")
	fmt.Println(meeting.ID)

	if !meeting.StartTime.IsZero() {
		_, _ = labelStyle.Print("Date:      ")
		fmt.Println(meeting.StartTime.Format("2 Jan 2006 15:04"))
	}

	if len(meeting.Attendees) > 0 {
		_, _ = labelStyle.Print("Attendees: ")
		fmt.Println(strings.Join(meeting.Attendees, ", "))
	}

	if meeting.Summary != "" {
		fmt.Println()
		_, _ = labelStyle.Println("─── Summary ───")
		fmt.Println()
		if err := RenderMarkdown(meeting.Summary); err != nil {
			fmt.Println(meeting.Summary)
		}
	}

	if meeting.Notes != "" {
		fmt.Println()
		_, _ = labelStyle.Println("─── Notes ───")
		fmt.Println()
		if err := RenderMarkdown(meeting.Notes); err != nil {
			fmt.Println(meeting.Notes)
		}
	}

	return nil
}

func PrintError(err error) {
	errStyle := color.New(color.FgRed, color.Bold)
	_, _ = errStyle.Fprint(os.Stderr, "Error: ")
	_, _ = fmt.Fprintln(os.Stderr, err.Error())
}

func PrintSuccess(message string) {
	successStyle := color.New(color.FgGreen)
	_, _ = successStyle.Print("✓ ")
	fmt.Println(message)
}

func PrintWarning(message string) {
	warnStyle := color.New(color.FgYellow)
	_, _ = warnStyle.Print("⚠ ")
	fmt.Println(message)
}

func PrintInfo(message string) {
	infoStyle := color.New(color.Faint)
	_, _ = infoStyle.Println(message)
}

type UserError struct {
	Message string
}

func (e *UserError) Error() string {
	return e.Message
}

func PrintTranscript(text string) error {
	var data struct {
		Title      string `json:"title"`
		Transcript string `json:"transcript"`
	}
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		// Not JSON, render as markdown
		return RenderMarkdown(text)
	}

	if data.Title != "" {
		titleStyle := color.New(color.Bold, color.FgWhite)
		_, _ = titleStyle.Println(data.Title)
		fmt.Println()
	}

	transcript := data.Transcript
	// Format speaker labels
	transcript = strings.ReplaceAll(transcript, "  Me:", "\n\n**Me:**")
	transcript = strings.ReplaceAll(transcript, "  Them:", "\n\n**Them:**")

	return RenderMarkdown(strings.TrimSpace(transcript))
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < 24*time.Hour && t.Day() == now.Day():
		return "today"
	case diff < 48*time.Hour && t.Day() == now.Add(-24*time.Hour).Day():
		return "yesterday"
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("2 Jan 2006")
	}
}


