package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

// readConfig reads and unmarshals the configuration file.
func ReadConfig(filePath string) ([]ExchangeConfig, error) {
	rawConfig, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var configs []ExchangeConfig
	if err := json.Unmarshal(rawConfig, &configs); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return configs, nil
}

// gracefulShutdown waits for a termination signal and closes all connections.
func GracefulShutdown(connections []*websocket.Conn, logger *log.Logger) {
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
