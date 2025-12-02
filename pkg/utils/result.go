package utils //nolint:revive //TODO: figure out a better name for this package

import (
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewToolResultText(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: message,
			},
		},
	}
}

func NewToolResultError(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: message,
			},
		},
		IsError: true,
	}
}

func NewToolResultErrorFromErr(message string, err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: message + ": " + err.Error(),
			},
		},
		IsError: true,
	}
}

func NewToolResultResource(message string, contents *mcp.ResourceContents) *mcp.CallToolResult {
	summary := message
	if metadata := formatResourceMetadata(contents); metadata != "" {
		summary = fmt.Sprintf("%s\n%s", message, metadata)
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: summary,
			},
		},
		IsError: false,
	}

	if contents == nil {
		return result
	}

	switch {
	case contents.Text != "":
		result.Content = append(result.Content, &mcp.TextContent{
			Text: contents.Text,
		})
	case len(contents.Blob) > 0 && strings.HasPrefix(contents.MIMEType, "image/"):
		result.Content = append(result.Content, &mcp.ImageContent{
			Data:     contents.Blob,
			MIMEType: contents.MIMEType,
		})
	default:
		// For non-text, non-image content we intentionally return metadata only.
	}

	return result
}

func formatResourceMetadata(contents *mcp.ResourceContents) string {
	if contents == nil {
		return ""
	}

	var parts []string

	if contents.URI != "" {
		parts = append(parts, fmt.Sprintf("URI: %s", contents.URI))
	}

	if contents.MIMEType != "" {
		parts = append(parts, fmt.Sprintf("MIME: %s", contents.MIMEType))
	}

	if len(contents.Blob) > 0 {
		parts = append(parts, fmt.Sprintf("Size: %d bytes", len(contents.Blob)))
	}

	return strings.Join(parts, " | ")
}
