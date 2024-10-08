package main

import (
	"fmt"
	"log"
	"net/http"
)

var messageChan chan string

func handleSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Get handshake from client")

		// prepare the headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// instantiate the channel
		messageChan = make(chan string)

		// close the channel after exit the function
		defer func() {
			close(messageChan)
			messageChan = nil
			log.Printf("client connection closed")
		}()

		// prepare the flusher
		flusher, _ := w.(http.Flusher)

		// trap the request under loop forever
		for {
			select {

			// message will be received here and printed
			case message := <-messageChan:
				fmt.Fprintf(w, "data: %s\n\n", message)
				flusher.Flush()

			// connection is closed and then defer will be executed
			case <-r.Context().Done():
				return
			}
		}
	}
}

func sendMessage(message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if messageChan != nil {
			log.Printf("print message to client")

			// send the message through available channnel
			messageChan <- message
		}
	}
}

func main() {
	http.HandleFunc("/handshake", handleSSE())

	http.HandleFunc("/sendmessage", sendMessage("hello client"))

	log.Fatal("HTTP server error: ", http.ListenAndServe(":8080", nil))
}
