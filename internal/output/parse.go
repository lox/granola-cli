package output

import (
	"encoding/xml"
	"regexp"
	"strings"
	"time"
)

// XML structures matching Granola MCP response format

type xmlMeetingsData struct {
	XMLName  xml.Name     `xml:"meetings_data"`
	From     string       `xml:"from,attr"`
	To       string       `xml:"to,attr"`
	Count    int          `xml:"count,attr"`
	Meetings []xmlMeeting `xml:"meeting"`
}

type xmlMeeting struct {
	ID           string `xml:"id,attr"`
	Title        string `xml:"title,attr"`
	Date         string `xml:"date,attr"`
	Participants string `xml:"known_participants"`
	Summary      string `xml:"summary"`
	Notes        string `xml:"notes"`
	PrivateNotes string `xml:"private_notes"`
}

var dateFormats = []string{
	"Jan 2, 2006 3:04 PM",
	"Jan 02, 2006 3:04 PM",
	"Jan 2, 2006 15:04",
	"2006-01-02T15:04:05Z",
	time.RFC3339,
}

func parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	for _, fmt := range dateFormats {
		if t, err := time.Parse(fmt, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// sanitizeXML fixes invalid XML from Granola's MCP responses.
// Granola includes unescaped <> in attribute values (e.g. title="A <> B"),
// <email@domain.com> in text content, and bare & in text, which break XML parsing.
func sanitizeXML(text string) string {
	// Escape < and > inside attribute values (between quotes after =)
	text = regexp.MustCompile(`(=")([^"]*?)(")`).ReplaceAllStringFunc(text, func(match string) string {
		prefix := match[:2] // ="
		suffix := match[len(match)-1:]
		inner := match[2 : len(match)-1]
		inner = strings.ReplaceAll(inner, "&", "&amp;")
		inner = strings.ReplaceAll(inner, "<", "&lt;")
		inner = strings.ReplaceAll(inner, ">", "&gt;")
		return prefix + inner + suffix
	})

	// Escape <email@domain> patterns in text content
	text = regexp.MustCompile(`<([a-zA-Z0-9_.+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,})>`).
		ReplaceAllString(text, "&lt;$1&gt;")

	// Escape <resource_id@resource.calendar.google.com> and similar long addresses
	text = regexp.MustCompile(`<([a-zA-Z0-9_]+@[a-zA-Z0-9._]+)>`).
		ReplaceAllString(text, "&lt;$1&gt;")

	// Escape bare & in text content (not already part of &amp; &lt; &gt; etc.)
	text = regexp.MustCompile(`&(amp;|lt;|gt;|quot;|apos;)`).
		ReplaceAllString(text, "\x00$1") // temporarily protect valid entities
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "\x00", "&") // restore valid entities

	return text
}

func ParseMeetingsList(text string) ([]Meeting, error) {
	text = sanitizeXML(text)

	var data xmlMeetingsData
	if err := xml.Unmarshal([]byte(text), &data); err != nil {
		return nil, &UserError{Message: "could not parse meetings response: " + err.Error()}
	}

	meetings := make([]Meeting, 0, len(data.Meetings))
	for _, m := range data.Meetings {
		meetings = append(meetings, Meeting{
			ID:        m.ID,
			Title:     m.Title,
			StartTime: parseDate(m.Date),
		})
	}
	return meetings, nil
}

func ParseMeetingDetail(text string) (*MeetingDetail, error) {
	text = sanitizeXML(text)

	var data xmlMeetingsData
	if err := xml.Unmarshal([]byte(text), &data); err != nil {
		return nil, &UserError{Message: "could not parse meeting detail: " + err.Error()}
	}

	if len(data.Meetings) == 0 {
		return nil, &UserError{Message: "no meeting found"}
	}

	m := data.Meetings[0]
	detail := &MeetingDetail{
		ID:        m.ID,
		Title:     m.Title,
		StartTime: parseDate(m.Date),
		Summary:   strings.TrimSpace(m.Summary),
		Notes:     strings.TrimSpace(m.Notes),
		Attendees: parseParticipants(m.Participants),
	}

	return detail, nil
}

func parseParticipants(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "Unknown" {
		return nil
	}

	parts := strings.Split(raw, ",")
	attendees := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		// Strip "(note creator)" and org info
		p = strings.ReplaceAll(p, "(note creator)", "")
		if idx := strings.Index(p, " from "); idx > 0 {
			p = p[:idx]
		}
		p = strings.TrimSpace(p)
		if p != "" {
			attendees = append(attendees, p)
		}
	}
	return attendees
}
