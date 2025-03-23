package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/txj-xyz/neoserve/file-server/internal/config"
)

// Embed represents a Discord embed structure.
type Embed struct {
	Title  string  `json:"title,omitempty"`
	URL    string  `json:"url,omitempty"`
	Color  int     `json:"color,omitempty"`
	Fields []Field `json:"fields,omitempty"`
	Footer *Footer `json:"footer,omitempty"`
	Image  *Image  `json:"image,omitempty"`
}

// Field represents a field in a Discord embed.
type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Footer represents the footer of a Discord embed.
type Footer struct {
	Text string `json:"text,omitempty"`
}

// Image represents an image in a Discord embed.
type Image struct {
	URL string `json:"url,omitempty"`
}

// WebhookPayload represents the full Discord webhook payload.
type WebhookPayload struct {
	Content     string  `json:"content,omitempty"`
	Embeds      []Embed `json:"embeds,omitempty"`
	Attachments []any   `json:"attachments,omitempty"` // Placeholder for future file support
}

// WebhookBuilder provides a builder pattern for creating webhook payloads.
type WebhookBuilder struct {
	payload WebhookPayload
}

// NewWebhook initializes a new WebhookBuilder.
func NewWebhook() *WebhookBuilder {
	return &WebhookBuilder{payload: WebhookPayload{}}
}

// SetContent sets the content of the webhook message.
func (wb *WebhookBuilder) SetContent(content string) *WebhookBuilder {
	wb.payload.Content = content
	return wb
}

// AddEmbed adds an embed to the webhook payload.
func (wb *WebhookBuilder) AddEmbed(embed Embed) *WebhookBuilder {
	wb.payload.Embeds = append(wb.payload.Embeds, embed)
	return wb
}

// Build constructs the final WebhookPayload.
func (wb *WebhookBuilder) Build() WebhookPayload {
	return wb.payload
}

// Send sends the webhook payload to the specified Discord webhook URL.
func (wb *WebhookBuilder) Send(webhookURL string) error {
	payloadBytes, err := json.Marshal(wb.payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("received error status code: %d", resp.StatusCode)
	}

	return nil
}

func SendWebhookLog(link string, fileName string) {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("cancelling webhook call: '%s'", err)
		return
	}

	if !cfg.Logging.Discord.Enabled {
		return
	}

	builder := NewWebhook()
	embed := Embed{
		Title: "**File Uploaded**",
		URL:   link,
		Fields: []Field{
			{
				Name:  "Name",
				Value: fmt.Sprintf("`%s`", fileName),
			},
			{
				Name:  "**Link**",
				Value: fmt.Sprintf("[Uncompressed Original](%s)", link),
			},
		},
		Footer: &Footer{
			Text: "Powered by Neoserve",
		},
		Image: &Image{
			URL: link,
		},
	}

	builder.SetContent("").AddEmbed(embed)
	err = builder.Send(cfg.Logging.Discord.WebhookURL)
	if err != nil {
		fmt.Printf("Failed to send discord webhook log: %s\n", err)
	}
}
