package binance

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/gorilla/websocket"
)

// gets executed as a goroutine
func HandleConnection(conn *websocket.Conn, exchange structs.ExchangeConfig) {
	if conn == nil {
		log.Println("Handle connection executed with no connection")
		return
	}

	// Parse streams from exchange config
	var streams []map[string]interface{}
	if err := json.Unmarshal(exchange.Streams, &streams); err != nil {
		log.Printf("Error unmarshalling streams for %s: %s\n", exchange.Name, err)
		return
	}

	// Send subscribe messages for each stream
	for _, stream := range streams {
		message, err := json.Marshal(stream)
		if err != nil {
			log.Printf("Error marshalling subscribe message for %s: %s\n", exchange.Name, err)
			continue
		}
		if err2 := conn.WriteMessage(websocket.TextMessage, message); err2 != nil {
			log.Printf("Error sending subscribe message for %s: %s\n", exchange.Name, err2)
		}
	}

	// Graceful shutdown handling
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt) // Cleanup signal handling

	done := make(chan struct{})

	messageQueue := make(chan []byte, 100)

	go func() {
		for message := range messageQueue {
			HandleMessage(message, exchange)
		}
	}()

	// Starts goroutine to start receiving incoming messages
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error receiving message from exchange %s: %s\n", exchange.Name, err)
				close(messageQueue) // Close message queue to allow worker to exit
				return
			}

			// Send message to the queue instead of spawning unlimited goroutines
			select {
			case messageQueue <- message:
			default:
				log.Println("Message queue full, dropping message")
			}
		}
	}()

	<-interrupt
	log.Println("Interrupt received, closing connection...")

	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Printf("Error closing WebSocket connection for %s: %s\n", exchange.Name, err)
	}

	conn.Close()
}

// Little helper function to process message types
func processMessageType(eventType string, message []byte, exchange string) error {
	// params handling

	switch eventType {
	case "24hrTicker":
		TickerData, err := HandleTickerMessage(message, true)
		if err != nil {
			return err
		}
		appendTickerBuffer(TickerData, exchange)
	case "trade":
		TradeData, err := HandleTradeMessage(message, true)
		if err != nil {
			return err
		}
		appendTradeBuffer(TradeData, exchange)
	default:
		log.Printf("Unhandled event type: %s | Exchange: %s\n", eventType, exchange)
	}

	return nil
}

// Check if message is wrapped or not
func HandleMessage(message []byte, exchange structs.ExchangeConfig) {
	var cMessage MessageCheck
	if err := json.Unmarshal(message, &cMessage); err != nil {
		log.Printf("Error parsing message: %s | Exchange: %s | Data: %s\n", err, exchange.Name, string(message)[:100])
		return
	}

	if cMessage.Stream == "" {
		// Unwrapped message
		var pMessage USMessageStruct
		if err := json.Unmarshal(message, &pMessage); err != nil {
			log.Printf("Error parsing US message: %s | Data: %s\n", err, string(message)[:100])
			return
		}

		if err := processMessageType(pMessage.EventType, message, "Binance_US"); err != nil {
			log.Printf("Error Processing Message Type") // fix this
			return
		}

	} else {
		// Wrapped message
		var pMessage GlobalMessageStruct
		if err := json.Unmarshal(message, &pMessage); err != nil {
			log.Printf("Error parsing global message: %s | Data: %s\n", err, string(message)[:100])
			return
		}

		processMessageType(pMessage.Data.EventType, message, "Binance_Global")
	}
}

// Handle Ticker Messages
func HandleTickerMessage(message []byte, wrapped bool) (structs.TickerData, error) {
	var pData structs.TickerData
	if wrapped {

	} else {

	}

	return pData, nil
}

// Handle Trade Messages
func HandleTradeMessage(message []byte, wrapped bool) (structs.TradeData, error) {
	var pData structs.TradeData

	if wrapped {

	} else {

	}

	return pData, nil
}

// Gracefully close the connection
func CloseConnection(conn *websocket.Conn) {
	if conn == nil {
		log.Println("CloseConnection called with nil connection")
		return
	}

	// Create a close message with a normal closure code.
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Normal closure")

	// Send the close message to the server.
	if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
		log.Printf("Error sending close message: %v", err)
		// Even if sending the close message fails, we'll proceed to close the connection.
	}

	// Optionally wait a moment to allow the close handshake to complete.
	// This gives the server a chance to acknowledge our close message.
	time.Sleep(1 * time.Second)

	// Close the underlying connection.
	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	} else {
		log.Println("Connection closed gracefully")
	}
}
