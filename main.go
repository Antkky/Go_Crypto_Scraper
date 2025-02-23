package main

import (
	"log"
	"os"
	"strings"

	"github.com/Antkky/go_crypto_scraper/handlers/binance"
	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/gorilla/websocket"
)

var logger = log.New(os.Stdout, "[CryptoScraper] ", log.LstdFlags|log.Lshortfile)

// establishConnections creates WebSocket connections to all the exchanges.
func establishConnections(configs []utils.ExchangeConfig) ([]*websocket.Conn, error) {
	var connections []*websocket.Conn

	for _, config := range configs {
		conn, _, err := websocket.DefaultDialer.Dial(config.URI, nil)
		if err != nil || conn == nil {
			logger.Printf("❌ Error connecting to exchange %s: %s", config.Name, err)
			continue
		} else {
			logger.Printf("✅ Connection established for %s", config.Name)
			connections = append(connections, conn)
			go handleExchangeConnection(config, conn)
		}
	}

	return connections, nil
}

// handleExchangeConnection routes connection handling based on the exchange.
func handleExchangeConnection(config utils.ExchangeConfig, conn *websocket.Conn) {
	switch {
	case strings.Contains(config.Name, "Binance"):
		binance.HandleConnection(conn, config, logger)
	case strings.Contains(config.Name, "Coinex"):
		//coinex.HandleConnection(conn, config)
	case strings.Contains(config.Name, "Bybit"):
		//bybit.HandleConnection(conn, config)
	case strings.Contains(config.Name, "Bitfinex"):
		//bitfinex.HandleConnection(conn, config)
	default:
		logger.Printf("⚠️ Unhandled exchange: %s", config.Name)
	}
}

func main() {
	// Read and parse configuration
	configs, err := utils.ReadConfig("config/streams2.json")
	if err != nil {
		logger.Fatalf("Error loading config: %s", err)
	}

	// Establish WebSocket connections
	connections, err := establishConnections(configs)
	if err != nil {
		logger.Fatalf("Error establishing connections: %s", err)
	}

	// Graceful shutdown handling
	utils.GracefulShutdown(connections, logger)
}
