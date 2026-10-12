package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type LoggedRequest struct {
	Timestamp string              `json:"timestamp"`
	Method    string              `json:"method"`
	URL       string              `json:"url"`
	Headers   map[string][]string `json:"headers"`
	RemoteAddr string             `json:"remote_addr"`
	Body      string              `json:"body"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are supported", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	entry := LoggedRequest{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Method:     r.Method,
		URL:        r.URL.String(),
		Headers:    r.Header,
		RemoteAddr: r.RemoteAddr,
		Body:       string(body),
	}

	encoder := json.NewEncoder(log.Writer())
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(entry); err != nil {
		log.Printf("Failed to encode request: %v", err)
		http.Error(w, "Failed to log request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK\n"))
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
