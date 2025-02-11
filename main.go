package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Antkky/go_crypto_scraper/binance"
	"github.com/Antkky/go_crypto_scraper/coinex"

	"github.com/gorilla/websocket"
)

type WebSocketConfig []ExchangeConfig

func gracefulShutdown(connections []*websocket.Conn) {
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

func connectExchange(exchange ExchangeConfig) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(exchange.URI, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func routeSubscribe(exchange ExchangeConfig, conn *websocket.Conn) error {
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

func routeResponse(conn *websocket.Conn, exchange ExchangeConfig) error {
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %s\n", exchange.Name, err)
			return err
		}
		if strings.HasPrefix(exchange.Name, "Binance") {
			binance.HandleMessage(message)
		} else if strings.HasPrefix(exchange.Name, "Coinex") {
			coinex.HandleMessage(message)
		} else {
			log.Println("Unhandled Exchange Response Type")
		}
	}
}

func main() {
	// Open & Parse Config File
	raw_config, err := os.ReadFile("config/streams.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s\n", err)
	}
	var config WebSocketConfig
	if err := json.Unmarshal(raw_config, &config); err != nil {
		log.Fatalf("Error parsing JSON: %s\n", err)
	}

	// Array to hold the connections
	connections := make([]*websocket.Conn, len(config))

	// Connect & Subscribe to Exchanges
	for i, exchange := range config {
		conn, err1 := connectExchange(exchange)
		if err1 != nil {
			log.Printf("Error connecting to exchange: %s\nError: %s\n", exchange.Name, err1)
			continue
		}
		if conn != nil {
			connections[i] = conn
			log.Println("Connected to: ", exchange.Name)
		}
		if err2 := routeSubscribe(exchange, conn); err2 != nil {
			log.Printf("Error subscribing to exchange: %s\nError: %s\n", exchange.Name, err2)
		}
		log.Println("Subscribed to: ", exchange.Name)

		// Remember to add a ping task to these connections (exchange.ping)
	}

	// Check for responses
	for i, connection := range connections {
		exchange := config[i]
		go routeResponse(connection, exchange)
		println("Listening on: ", exchange.Name)
	}

	go gracefulShutdown(connections)
	select {}
}
