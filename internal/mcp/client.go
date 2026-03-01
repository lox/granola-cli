package mcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	DefaultEndpoint = "https://mcp.granola.ai/mcp"
)

type Client struct {
	mcpClient  *client.Client
	tokenStore *FileTokenStore
}

type ClientOption func(*clientConfig)

type clientConfig struct {
	endpoint    string
	accessToken string
}

func WithEndpoint(endpoint string) ClientOption {
	return func(c *clientConfig) {
		c.endpoint = endpoint
	}
}

func WithAccessToken(token string) ClientOption {
	return func(c *clientConfig) {
		c.accessToken = token
	}
}

func NewClient(opts ...ClientOption) (*Client, error) {
	cfg := &clientConfig{
		endpoint: DefaultEndpoint,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	tokenStore, err := NewFileTokenStore()
	if err != nil {
		return nil, fmt.Errorf("create token store: %w", err)
	}

	var store transport.TokenStore = tokenStore
	if cfg.accessToken != "" {
		store = &staticTokenStore{token: cfg.accessToken}
	}

	oauthConfig := transport.OAuthConfig{
		TokenStore:  store,
		PKCEEnabled: true,
	}

	trans, err := transport.NewStreamableHTTP(
		cfg.endpoint,
		transport.WithHTTPOAuth(oauthConfig),
	)
	if err != nil {
		return nil, fmt.Errorf("create transport: %w", err)
	}

	return &Client{
		mcpClient:  client.NewClient(trans),
		tokenStore: tokenStore,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	if err := c.mcpClient.Start(ctx); err != nil {
		if client.IsOAuthAuthorizationRequiredError(err) {
			return &AuthRequiredError{
				Handler: client.GetOAuthHandler(err),
			}
		}
		return err
	}

	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "granola-cli",
		Version: "0.1.0",
	}

	_, err := c.mcpClient.Initialize(ctx, initReq)
	if err != nil {
		if client.IsOAuthAuthorizationRequiredError(err) {
			return &AuthRequiredError{
				Handler: client.GetOAuthHandler(err),
			}
		}
		return fmt.Errorf("initialize: %w", err)
	}

	return nil
}

func (c *Client) Close() error {
	return c.mcpClient.Close()
}

func (c *Client) TokenStore() *FileTokenStore {
	return c.tokenStore
}

func (c *Client) GetOAuthHandler() *transport.OAuthHandler {
	trans := c.mcpClient.GetTransport()
	if st, ok := trans.(*transport.StreamableHTTP); ok {
		return st.GetOAuthHandler()
	}
	return nil
}

type AuthRequiredError struct {
	Handler *transport.OAuthHandler
}

func (e *AuthRequiredError) Error() string {
	return "authentication required - run 'granola auth login'"
}

func IsAuthRequired(err error) bool {
	var authErr *AuthRequiredError
	return errors.As(err, &authErr)
}

func GetOAuthHandler(err error) *transport.OAuthHandler {
	var authErr *AuthRequiredError
	if errors.As(err, &authErr) {
		return authErr.Handler
	}
	return nil
}

func (c *Client) CallTool(ctx context.Context, name string, args map[string]any) (*mcp.CallToolResult, error) {
	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args

	return c.mcpClient.CallTool(ctx, req)
}

func (c *Client) CallToolText(ctx context.Context, name string, args map[string]any) (string, error) {
	result, err := c.CallTool(ctx, name, args)
	if err != nil {
		return "", err
	}
	if err := checkToolError(result); err != nil {
		return "", err
	}
	return extractText(result), nil
}

func (c *Client) ListTools(ctx context.Context) ([]mcp.Tool, error) {
	resp, err := c.mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Tools, nil
}

// staticTokenStore provides a token from a fixed string (for CI/env var usage)
type staticTokenStore struct {
	token string
}

func (s *staticTokenStore) GetToken(ctx context.Context) (*transport.Token, error) {
	return &transport.Token{
		AccessToken: s.token,
		TokenType:   "Bearer",
	}, nil
}

func (s *staticTokenStore) SaveToken(ctx context.Context, token *transport.Token) error {
	return nil
}

func checkToolError(result *mcp.CallToolResult) error {
	if result == nil || !result.IsError {
		return nil
	}
	msg := extractText(result)
	if msg == "" {
		msg = "tool call failed"
	}
	return fmt.Errorf("granola API error: %s", msg)
}

func extractText(result *mcp.CallToolResult) string {
	if result == nil {
		return ""
	}
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			return textContent.Text
		}
	}
	return ""
}
