package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

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

func revised_main() {
	/******** Open & Parse Config File *********/
	// open
	rawConfig, err := os.ReadFile("config/streams.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s\n", err)
		return
	}
	// declare
	var configs []ExchangeConfig
	// parse
	if err = json.Unmarshal(rawConfig, &configs); err != nil {
		log.Fatalf("Error parsing JSON: %s\n", err)
		return
	}
	/********** Establish Connections **********/
	connections := make([]*websocket.Conn, len(configs))
	for i, config := range configs {
		// Connect
		// refactor the function?
		conn, err := connectExchange(config)
		if err != nil {
			log.Printf("Error connecting to exchange: %s\nError: %s\n", config.Name, err)
			continue
		}
		if conn != nil {
			connections[i] = conn
		} else {
			log.Printf("Error, connection for exchange: %s is nil.", config.Name)
		}
		// Subscribe
		// refactor the function?
		if err2 := routeSubscribe(config, conn); err2 != nil {
			log.Printf("Error subscribing to exchange: %s\nError: %s\n", config.Name, err2)
		}
		// Ping
		// function not implemented yet
		log.Printf("Connection Established for: %s", config.Name)
	}
	/******** Launch Connection Handlers *******/

	/*********** Graceful Shutdown ************/
	go gracefulShutdown(connections)
	select {} // fix this
}
