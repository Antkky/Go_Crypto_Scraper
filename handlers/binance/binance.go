package binance

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/gorilla/websocket"
)

// HandleConnection manages the lifecycle of a WebSocket connection to Binance.
func HandleConnection(conn *websocket.Conn, exchange utils.ExchangeConfig) {
	if conn == nil {
		log.Println("Connection is nil, exiting HandleConnection.")
		return
	}

	streams, err := parseStreams(exchange.Streams)
	if err != nil {
		log.Printf("Failed to parse streams for %s: %s", exchange.Name, err)
		return
	}

	// Subscribe to the streams
	subscribeToStreams(conn, streams, exchange.Name)

	// Setup graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	messageQueue := make(chan []byte, 100)
	done := make(chan struct{})

	// Concurrent message handling
	go ConsumeMessages(messageQueue, done, exchange)
	go ReceiveMessages(conn, messageQueue, done, exchange)

	// Wait for interrupt
	<-interrupt
	log.Println("Interrupt received, closing connection...")

	// Cleanly close connection
	closeConnection(conn, exchange.Name)
}

func parseStreams(streamsData json.RawMessage) ([]map[string]interface{}, error) {
	var streams []map[string]interface{}
	if err := json.Unmarshal(streamsData, &streams); err != nil {
		return nil, err
	}
	return streams, nil
}

func subscribeToStreams(conn *websocket.Conn, streams []map[string]interface{}, exchangeName string) {
	for _, stream := range streams {
		if err := Subscribe(conn, stream, exchangeName); err != nil {
			log.Printf("Error subscribing to stream for %s: %s", exchangeName, err)
		}
	}
}

func Subscribe(conn *websocket.Conn, stream map[string]interface{}, exchangeName string) error {
	message, err := json.Marshal(stream)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}
	return nil
}

func ConsumeMessages(messageQueue chan []byte, done chan struct{}, exchange utils.ExchangeConfig) {
	defer close(done)
	for message := range messageQueue {
		if err := HandleMessage(message, exchange); err != nil {
			log.Printf("Error handling message for %s: %v", exchange.Name, err)
		}
	}
}

func ReceiveMessages(conn *websocket.Conn, messageQueue chan []byte, done chan struct{}, exchange utils.ExchangeConfig) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", exchange.Name, err)
			close(messageQueue)
			return
		}
		select {
		case messageQueue <- message:
		default:
			log.Printf("Message queue full, dropping message for %s", exchange.Name)
		}
	}
}

func ProcessMessageType(message []byte, tickerData *utils.TickerDataStruct, tradeData *utils.TradeDataStruct) error {
	var pMessage GlobalMessageStruct
	var eventType string

	if pMessage.Data.EventType == "" {
		eventType = pMessage.EventType
	} else {
		eventType = pMessage.Data.EventType
	}

	switch eventType {
	case "24hrTicker":
		return HandleTickerMessage(message, tickerData)
	case "trade":
		return HandleTradeMessage(message, tradeData)
	default:
		return errors.New("unhandled event type: " + eventType)
	}
}

func HandleTickerMessage(message []byte, tickerData *utils.TickerDataStruct) error {
	// Add logic here to handle ticker message
	// e.g., Unmarshal and process the data, based on whether it's wrapped or not
	return nil
}

func HandleTradeMessage(message []byte, tradeData *utils.TradeDataStruct) error {
	// Add logic here to handle trade message
	// e.g., Unmarshal and process the data, based on whether it's wrapped or not
	return nil
}

func HandleMessage(message []byte, exchange utils.ExchangeConfig) error {
	var (
		tickerData utils.TickerDataStruct
		tradeData  utils.TradeDataStruct
	)

	if err := ProcessMessageType(message, &tickerData, &tradeData); err != nil {
		return err
	}

	return nil
}

func closeConnection(conn *websocket.Conn, exchangeName string) {
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Normal closure")
	if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
		log.Printf("Error sending close message for %s: %v", exchangeName, err)
	}

	time.Sleep(time.Second) // Give some time for the close to be processed

	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection for %s: %v", exchangeName, err)
	} else {
		log.Printf("Connection for %s closed gracefully", exchangeName)
	}
}
