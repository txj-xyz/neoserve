package main

import (
	"fmt"

	"github.com/txj-xyz/neoserve/webhook"
)

// Main function for testing the builder.
func main() {

	builder := webhook.NewWebhook()

	webhookURL := "https://discord.com/api/webhooks/1308482375243661372/37RPD53RCrqiPgrWRgDzxR59oAM5br2WC5rvNFfntQ0oSzHiW9r0UsBg4_Dof5h8XRsn"


	embed := webhook.Embed{
		Title: "**File Uploaded**",
		URL:   "https://l.txxj.xyz/v1/files/7cc65fa5c2e30bad90aa90eee382e72e.png",
		Color: 720640,
		Fields: []webhook.Field{
			{
				Name:  "Name",
				Value: "`7cc65fa5c2e30bad90aa90eee382e72e`",
			},
			{
				Name:  "**Link**",
				Value: "[Uncompressed Original](https://l.txxj.xyz/v1/files/7cc65fa5c2e30bad90aa90eee382e72e.png)",
			},
		},
		Footer: &webhook.Footer{
			Text: "Powered by Neoserve",
		},
		Image: &webhook.Image{
			URL: "https://l.txxj.xyz/v1/files/7cc65fa5c2e30bad90aa90eee382e72e.png",
		},
	}


	err := builder.SetContent("").AddEmbed(embed).Send(webhookURL)
	if err != nil {
		fmt.Println("Failed to send webhook")
		return
	}



	// Send the embed to the webhook.

	fmt.Println("Webhook sent successfully!")
}