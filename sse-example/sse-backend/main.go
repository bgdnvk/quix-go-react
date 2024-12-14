package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type DataItem struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

var globalClients = make(map[chan string]bool)

func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Immediately return data
	items := []DataItem{
		{ID: 1, Value: "A"},
		{ID: 2, Value: "B"},
		{ID: 3, Value: "C"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)

	// Launch validations in background
	go validateData(items)
}

func validateData(items []DataItem) {
	// Simulate parallel validations
	for _, item := range items {
		go func(it DataItem) {
			time.Sleep(time.Duration(it.ID) * time.Second) // simulate delay

			// Validation logic: say value can't be "B"
			valid := it.Value != "B"
			result := fmt.Sprintf("id:%d, valid:%t", it.ID, valid)

			// Broadcast via SSE
			broadcastSSE(result)
		}(item)
	}
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)
	globalClients[messageChan] = true
	defer func() {
		delete(globalClients, messageChan)
		close(messageChan)
	}()

	// Keep connection open
	for {
		select {
		case msg := <-messageChan:
			// Write SSE event
			fmt.Fprintf(w, "event: validationResult\n")
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func broadcastSSE(msg string) {
	for clientChan := range globalClients {
		clientChan <- msg
	}
}

func main() {
	http.HandleFunc("/data", dataHandler)
	http.HandleFunc("/events", eventsHandler)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
