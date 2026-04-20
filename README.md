# Telegram Alert Dispatcher 🚨

A lightning-fast, standalone microservice built in **Go** that receives HTTP POST requests from your other services and asynchronously broadcasts beautiful, HTML-formatted alert messages to configured Telegram chats.

Perfect for centralizing logging, error tracking, and critical system notifications (like new transactions, server errors, or user registrations) across a distributed microservice architecture.

## 🚀 Key Features
- **Asynchronous Broadcasting:** Uses Go routines to dispatch messages instantly without blocking the HTTP response.
- **Dynamic HTML Formatting:** Automatically styles the Telegram message with appropriate emojis based on the alert level (`INFO`, `WARNING`, `ERROR`, `SUCCESS`).
- **Multi-Chat Support:** Broadcasts to one or multiple Telegram IDs/Groups simultaneously.
- **Zero External Dependencies:** Built entirely using Go's standard library (`net/http`, `encoding/json`).

## ⚙️ Setup & Installation

1. Clone the repository:
   ```bash
   git clone [https://github.com/card0re/tg-alert-dispatcher.git](https://github.com/card0re/tg-alert-dispatcher.git)
   cd tg-alert-dispatcher