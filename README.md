# Telegram Alert Dispatcher 🚨

A lightning-fast, standalone microservice built in **Go** that receives HTTP POST requests from your other services and asynchronously broadcasts beautiful, HTML-formatted alert messages to configured Telegram chats.

Perfect for centralizing logging, error tracking, and critical system notifications (like new transactions, server errors, or user registrations) across a distributed microservice architecture.

## 🚀 Key Features
- **Asynchronous Broadcasting:** Uses Go routines to dispatch messages instantly without blocking the HTTP response.
- **Dynamic HTML Formatting:** Automatically styles the Telegram message with appropriate emojis based on the alert level (`INFO`, `WARNING`, `ERROR`, `SUCCESS`).
- **Multi-Chat Support:** Broadcasts to one or multiple Telegram IDs/Groups simultaneously.
- **Zero External Dependencies:** Built entirely using Go's standard library (`net/http`, `encoding/json`).

## ⚙️ Run

```bash
git clone https://github.com/card0re/tg-alert-dispatcher.git
cd tg-alert-dispatcher
export TELEGRAM_BOT_TOKEN=<bot token from @BotFather>
export TELEGRAM_CHAT_IDS=123456789,-100987654321   # comma-separated
export PORT=8082                                     # optional, default 8082
go run .
```

## 📡 API

```bash
curl -X POST http://localhost:8082/api/v1/alerts/send \
  -H "Content-Type: application/json" \
  -d '{"level":"ERROR","source":"PaymentGateway","message":"Payment failed","details":{"order_id":"42"}}'
```
