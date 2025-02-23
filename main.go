package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Antkky/go_crypto_scraper/handlers/binance"
	"github.com/Antkky/go_crypto_scraper/handlers/coinex"
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
		coinex.HandleConnection(conn, config, logger)
	case strings.Contains(config.Name, "Bybit"):
		//bybit.HandleConnection(conn, config)
	case strings.Contains(config.Name, "Bitfinex"):
		//bitfinex.HandleConnection(conn, config)
	default:
		logger.Printf("⚠️ Unhandled exchange: %s", config.Name)
	}
}

// gracefulShutdown waits for a termination signal and closes all connections.
func GracefulShutdown(connections []*websocket.Conn, configs []utils.ExchangeConfig, logger *log.Logger) {
	// Wait for interrupt signal to gracefully shutdown the application.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	logger.Println("⏳ Shutting down...")

	for i, conn := range connections {
		if conn != nil {
			switch {
			case strings.Contains(configs[i].Name, "Binance"):
				binance.CloseConnection(conn, configs[i].Name, logger)
			case strings.Contains(configs[i].Name, "Coinex"):
				coinex.CloseConnection(conn, configs[i].Name, logger)
			case strings.Contains(configs[i].Name, "Bybit"):
				//bybit.HandleConnection(conn, config)
			case strings.Contains(configs[i].Name, "Bitfinex"):
				//bitfinex.HandleConnection(conn, config)
			default:
				logger.Printf("⚠️ Unhandled exchange: %s", configs[i].Name)
			}
		}
	}

	logger.Println("✅ Cleanup complete. Exiting.")
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
	GracefulShutdown(connections, configs, logger)
}
