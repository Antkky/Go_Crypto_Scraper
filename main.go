package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Antkky/go_crypto_scraper/handlers/binance"
	"github.com/Antkky/go_crypto_scraper/handlers/bitfinex"
	"github.com/Antkky/go_crypto_scraper/handlers/bybit"
	"github.com/Antkky/go_crypto_scraper/handlers/coinex"
	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/gorilla/websocket"
)

// Error and logging handling improvements.
var logger = log.New(os.Stdout, "[CryptoScraper] ", log.LstdFlags|log.Lshortfile)

// readConfig reads and unmarshals the configuration file.
func readConfig(filePath string) ([]utils.ExchangeConfig, error) {
	rawConfig, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configs []utils.ExchangeConfig
	if err := json.Unmarshal(rawConfig, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return configs, nil
}

// establishConnections creates WebSocket connections to all the exchanges.
func establishConnections(configs []utils.ExchangeConfig) ([]*websocket.Conn, error) {
	var connections []*websocket.Conn

	for _, config := range configs {
		conn, _, err := websocket.DefaultDialer.Dial(config.URI, nil)
		if err != nil {
			logger.Printf("❌ Error connecting to exchange %s: %s", config.Name, err)
			continue
		}
		if conn == nil {
			logger.Printf("⚠️ Connection for exchange %s is nil.", config.Name)
			continue
		}

		logger.Printf("✅ Connection established for %s", config.Name)
		connections = append(connections, conn)

		// Handle connection in separate goroutines based on exchange
		go handleExchangeConnection(config, conn)
	}

	return connections, nil
}

// handleExchangeConnection routes connection handling based on the exchange.
func handleExchangeConnection(config utils.ExchangeConfig, conn *websocket.Conn) {
	switch {
	case strings.Contains(config.Name, "Binance"):
		binance.HandleConnection(conn, config)
	case strings.Contains(config.Name, "Coinex"):
		coinex.HandleConnection(conn, config)
	case strings.Contains(config.Name, "Bybit"):
		bybit.HandleConnection(conn, config)
	case strings.Contains(config.Name, "Bitfinex"):
		bitfinex.HandleConnection(conn, config)
	default:
		logger.Printf("⚠️ Unhandled exchange: %s", config.Name)
	}
}

// gracefulShutdown waits for a termination signal and closes all connections.
func gracefulShutdown(connections []*websocket.Conn) {
	// Wait for interrupt signal to gracefully shutdown the application.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	logger.Println("⏳ Shutting down...")

	for _, conn := range connections {
		if conn != nil {
			// Attempt graceful WebSocket closure.
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Shutdown")); err != nil {
				logger.Printf("⚠️ Failed to send close message to WebSocket connection: %s", err)
			}
			conn.Close()
		}
	}

	logger.Println("✅ Cleanup complete. Exiting.")
}

func main() {
	// Read and parse configuration
	configs, err := readConfig("config/streams.json")
	if err != nil {
		logger.Fatalf("Error loading config: %s", err)
	}

	// Establish WebSocket connections
	connections, err := establishConnections(configs)
	if err != nil {
		logger.Fatalf("Error establishing connections: %s", err)
	}

	// Graceful shutdown handling
	gracefulShutdown(connections)
}
