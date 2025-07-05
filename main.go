package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
)

var client pusher.Client

func decodeBase64(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatalf("Failed to decode base64: %v", err)
	}
	return string(decoded)
}

func initPusher() {
	_ = godotenv.Load()

	appID := os.Getenv("APP_ID")
	key := decodeBase64(os.Getenv("BASE64_KEY"))
	secret := decodeBase64(os.Getenv("BASE64_SECRET"))
	scheme := "https"
	if os.Getenv("USE_TLS") == "false" {
		scheme = "http"
	}
	host := os.Getenv("SOKETI_HOST")
	fullHost := scheme + "://" + host

	client = pusher.Client{
		AppID:  appID,
		Key:    key,
		Secret: secret,
		Host:   fullHost,
		Secure: scheme == "https",
	}
}

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
	Room    string `json:"room"`
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := client.Trigger(msg.Room, "message", msg)
	if err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		log.Println("Trigger error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	// Read raw body as []byte
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read body", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	// Use new AuthorizePrivateChannel
	response, err := client.AuthorizePrivateChannel(body)
	if err != nil {
		http.Error(w, "Auth failed", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
}

func main() {
	initPusher()

	http.HandleFunc("/send", handleSendMessage)
	http.HandleFunc("/pusher/auth", handleAuth)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
