package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	botToken string
	chatIDs  []string
)

// --- STRUCTURES ---

type AlertRequest struct {
	Level   string            `json:"level"`   // e.g., "INFO", "WARNING", "ERROR", "SUCCESS"
	Source  string            `json:"source"`  // e.g., "PaymentGateway", "AuthService"
	Message string            `json:"message"` // Main alert text
	Details map[string]string `json:"details"` // Optional key-value metadata
}

// --- TELEGRAM LOGIC ---

// sendTelegramMessage sends a formatted HTML message to a specific chat asynchronously
func sendTelegramMessage(chatID, text string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	reqBody, _ := json.Marshal(map[string]string{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Printf("❌ Failed to send alert to %s: %v", chatID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️ Telegram API returned status %d for chat %s", resp.StatusCode, chatID)
	}
}

// formatAlertText creates a beautiful HTML template for the Telegram message
func formatAlertText(req AlertRequest) string {
	// Pick emoji based on level
	emoji := "ℹ️"
	switch strings.ToUpper(req.Level) {
	case "ERROR", "CRITICAL":
		emoji = "🚨"
	case "WARNING":
		emoji = "⚠️"
	case "SUCCESS":
		emoji = "✅"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s <b>[%s] ALERT</b>\n\n", emoji, strings.ToUpper(req.Level)))
	sb.WriteString(fmt.Sprintf("🏢 <b>Source:</b> <code>%s</code>\n", req.Source))
	sb.WriteString(fmt.Sprintf("📝 <b>Message:</b> %s\n", req.Message))

	if len(req.Details) > 0 {
		sb.WriteString("\n🔍 <b>Details:</b>\n")
		for k, v := range req.Details {
			sb.WriteString(fmt.Sprintf("├ <i>%s</i>: <code>%s</code>\n", k, v))
		}
	}

	sb.WriteString(fmt.Sprintf("\n🕒 <i>%s</i>", time.Now().Format(time.RFC1123)))
	return sb.String()
}

// --- HTTP HANDLERS ---

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AlertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	if req.Message == "" || req.Source == "" {
		http.Error(w, `{"error": "Missing required fields: 'source' or 'message'"}`, http.StatusBadRequest)
		return
	}

	// Format the message
	htmlText := formatAlertText(req)

	// Dispatch asynchronously to all configured chats
	for _, chatID := range chatIDs {
		go sendTelegramMessage(strings.TrimSpace(chatID), htmlText)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Alert dispatched to Telegram",
	})
}

func main() {
	botToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	chatsEnv := os.Getenv("TELEGRAM_CHAT_IDS")
	port := os.Getenv("PORT")

	if botToken == "" || chatsEnv == "" {
		log.Fatal("Environment variables TELEGRAM_BOT_TOKEN and TELEGRAM_CHAT_IDS are required")
	}

	// Split comma-separated chat IDs
	chatIDs = strings.Split(chatsEnv, ",")

	if port == "" {
		port = "8082"
	}

	http.HandleFunc("/api/v1/alerts/send", alertHandler)

	log.Printf("🚀 Telegram Alert Dispatcher running on port %s", port)
	log.Printf("📡 Configured to broadcast to %d chat(s)", len(chatIDs))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
