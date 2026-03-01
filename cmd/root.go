package cmd

type Context struct {
	JSON  bool
	Token string
}

type CLI struct {
	Token string `help:"Access token (skips OAuth)" env:"GRANOLA_ACCESS_TOKEN" hidden:""`

	Auth     AuthCmd     `cmd:"" help:"Authentication commands"`
	Meetings MeetingsCmd `cmd:"" help:"Meeting commands"`
	Query    QueryCmd    `cmd:"" help:"Query meetings using natural language"`
	Tools    ToolsCmd    `cmd:"" help:"List available MCP tools"`
	Version  VersionCmd  `cmd:"" help:"Show version"`
}

type VersionCmd struct {
	Version string `kong:"hidden,default='${version}'"`
}

func (c *VersionCmd) Run(ctx *Context) error {
	println("granola version " + c.Version)
	return nil
}
