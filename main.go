package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSockets.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// handleConnections handles incoming WebSocket requests from clients.
func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer ws.Close()

	log.Println("Client connected")

	for {
		messageType, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		flacData, err := convertWAVtoFLAC(msg)
		if err != nil {
			log.Println("Error converting WAV to FLAC:", err)
			continue
		}

		err = ws.WriteMessage(messageType, flacData)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}

// convertWAVtoFLAC converts WAV byte data to FLAC byte data using FFmpeg.
func convertWAVtoFLAC(wavData []byte) ([]byte, error) {
	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "flac", "pipe:1")
	cmd.Stdin = bytes.NewReader(wavData)

	var flacBuffer bytes.Buffer
	cmd.Stdout = &flacBuffer

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running FFmpeg: %v", err)
	}

	return flacBuffer.Bytes(), nil
}

func main() {
	http.HandleFunc("/stream", handleConnections)
	log.Println("WebSocket server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
