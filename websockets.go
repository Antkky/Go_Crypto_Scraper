package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/gorilla/websocket"
)

func RouteSubscribe(exchange structs.ExchangeConfig, conn *websocket.Conn) error {
	// Loop Through Streams
	var streams []map[string]interface{}
	if err := json.Unmarshal(exchange.Streams, &streams); err != nil {
		log.Printf("Error unmarshalling streams for %s: %s\n", exchange.Name, err)
		return err
	}
	// Loop Through Streams
	for _, stream := range streams {
		message, err := json.Marshal(stream)
		if err != nil {
			log.Printf("Error marshalling subscribe message for: %s, %s\n", exchange.Name, err)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error sending subscribe message for: %s, %s\n", exchange.Name, err)
		}
	}
	return nil
}

func GracefulShutdown(connections []*websocket.Conn) {
	// Create a channel to listen for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	<-sigChan

	// Close all connections gracefully
	log.Println("Shutting down gracefully...")
	for _, conn := range connections {
		if conn != nil {
			err := conn.Close()
			if err != nil {
				log.Printf("Error closing WebSocket connection: %s\n", err)
			}
		}
	}
	log.Println("All connections closed.")
}
