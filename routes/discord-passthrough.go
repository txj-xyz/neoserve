package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi"
)

func DiscordWebhookPassthrough(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	token := chi.URLParam(r, "webhookToken")
	discordWebhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", channelID, token)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	
	discordRequest, err := http.NewRequest("POST", discordWebhookURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	for header, values := range r.Header {
		for _, value := range values {
			discordRequest.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(discordRequest)
	if err != nil {
		http.Error(w, "Failed to forward request to Discord", http.StatusInternalServerError)
		return 
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		fmt.Printf("Failed to write response to client: %v", err)
	}

	// Log out the request
	// link := fmt.Sprintf("%s/v1/files/%s", cfg.Server.GenerateURL(), ".")
	// webhook.SendWebhookLog(link, rnd)
}