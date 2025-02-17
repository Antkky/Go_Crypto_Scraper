package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Antkky/go_crypto_scraper/binance"
	"github.com/Antkky/go_crypto_scraper/bitfinex"
	"github.com/Antkky/go_crypto_scraper/bybit"
	"github.com/Antkky/go_crypto_scraper/coinex"
	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/gorilla/websocket"
)

// test this i think
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

// test this
func ConnectExchange(exchange structs.ExchangeConfig) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(exchange.URI, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// test this
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

func main() {
	/******** Open & Parse Config File *********/
	rawConfig, err := os.ReadFile("config/streams.json")
	if err != nil {
		log.Fatalf("Error reading JSON file: %s\n", err)
	}

	var configs []structs.ExchangeConfig
	if err = json.Unmarshal(rawConfig, &configs); err != nil {
		log.Fatalf("Error parsing JSON: %s\n", err)
	}

	/********** Establish Connections **********/
	connections := make([]*websocket.Conn, 0) // Use slice instead of fixed-size array
	for _, config := range configs {
		conn, err := ConnectExchange(config)
		if err != nil {
			log.Printf("❌ Error connecting to exchange %s: %s\n", config.Name, err)
			continue
		}
		if conn == nil {
			log.Printf("⚠️ Connection for exchange %s is nil.\n", config.Name)
			continue
		}
		// Subscribe to streams
		if err := RouteSubscribe(config, conn); err != nil {
			log.Printf("⚠️ Error subscribing to exchange %s: %s\n", config.Name, err)
			conn.Close() // Close connection if subscription fails
			continue
		}

		log.Printf("✅ Connection Established for: %s", config.Name)
		connections = append(connections, conn)

		// Launch Connection Handler
		switch {
		case strings.Contains(config.Name, "Binance"):
			go binance.HandleConnection(conn, config)
		case strings.Contains(config.Name, "Coinex"):
			go coinex.HandleConnection(conn, config)
		case strings.Contains(config.Name, "Bybit"):
			go bybit.HandleConnection(conn, config)
		case strings.Contains(config.Name, "Bitfinex"):
			go bitfinex.HandleConnection(conn, config)
		default:
			log.Println("⚠️ Unhandled Exchange:", config.Name)
		}
	}

	/*********** Graceful Shutdown ************/
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan // Wait for termination signal
	log.Println("⏳ Shutting down...")

	// Close all WebSocket connections
	for _, conn := range connections {
		if conn != nil {
			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Shutdown"))
			conn.Close()
		}
	}

	log.Println("✅ Cleanup complete. Exiting.")
}
