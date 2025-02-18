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
		conn, _, err := websocket.DefaultDialer.Dial(config.URI, nil)
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
