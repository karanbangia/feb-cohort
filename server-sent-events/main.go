package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-faker/faker/v4"
)

var messageChan = make(chan []byte, 100)

func main() {
	go func() {
		for {
			name := faker.FirstName() + " " + faker.LastName()
			select {
			case messageChan <- []byte(fmt.Sprintf(`{"name": "%s"}`, name)):
			default:
				// Drop message if channel is full to avoid blocking
			}
			time.Sleep(time.Second) // Prevent excessive message generation
		}
	}()
	http.HandleFunc("/events", sseHandler)
	http.ListenAndServe(":8080", nil)
}

func sseHandler(rw http.ResponseWriter, req *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("X-Accel-Buffering", "no")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

	// keeping the connection alive with keep-alive protocol
	keepAliveTickler := time.NewTicker(15 * time.Second)
	defer keepAliveTickler.Stop()
	keepAliveMsg := ":keepalive\n"

	notify := req.Context().Done()

	for {
		select {
		// receiving a message from the Kafka channel.
		case kafkaEvent := <-messageChan:
			// Write to the ResponseWriter in SSE compatible format
			fmt.Fprintf(rw, "data: %s\n\n", kafkaEvent)
			flusher.Flush()
		case <-keepAliveTickler.C:
			fmt.Fprintf(rw, keepAliveMsg)
			flusher.Flush()
		case <-notify:
			fmt.Println("Client disconnected")
			return // Exit function to free resources
		}
	}

}
