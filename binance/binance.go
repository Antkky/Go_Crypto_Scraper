package binance

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/gorilla/websocket"
)

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
	defer signal.Stop(interrupt)
	done := make(chan struct{})
	messageQueue := make(chan []byte, 100)

	go func() {
		defer close(done)
		for message := range messageQueue {
			HandleMessage(message, exchange)
		}
	}()
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error receiving message from exchange %s: %s\n", exchange.Name, err)
				close(messageQueue)
				return
			}
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

func processMessageType(eventType string, message []byte, wrapped bool) (structs.TickerData, structs.TradeData, error) {
	if eventType == "" || len(message) == 0 {
		return structs.TickerData{}, structs.TradeData{}, errors.New("invalid parameters")
	}
	switch eventType {
	case "24hrTicker":
		TickerData, err := HandleTickerMessage(message, wrapped)
		if err != nil {
			return structs.TickerData{}, structs.TradeData{}, err
		}
		return TickerData, structs.TradeData{}, nil
	case "trade":
		TradeData, err := HandleTradeMessage(message, wrapped)
		if err != nil {
			return structs.TickerData{}, structs.TradeData{}, err
		}
		return structs.TickerData{}, TradeData, nil
	default:
		return structs.TickerData{}, structs.TradeData{}, errors.New("unhandled event type")
	}
}

func HandleMessage(message []byte, exchange structs.ExchangeConfig) error {
	var cMessage MessageCheck
	if err := json.Unmarshal(message, &cMessage); err != nil {
		log.Printf("Error parsing message: %s | Exchange: %s | Data: %s\n", err, exchange.Name, string(message)[:100])
		return err
	}
	if cMessage.Stream == "" {
		var pMessage USMessageStruct
		if err := json.Unmarshal(message, &pMessage); err != nil {
			log.Printf("Error parsing US message: %s | Data: %s\n", err, string(message)[:100])
			return err
		}
		tkrData, trdData, err := processMessageType(pMessage.EventType, message, false)
		if err != nil {
			log.Printf("Error Processing Message Type %s from exchange %s.\nError Code: %s", pMessage.EventType, exchange.Name, err)
			return err
		}
		if tkrData != (structs.TickerData{}) {

		} else if trdData != (structs.TradeData{}) {

		}
	} else if len(cMessage.Stream) >= 1 {
		var pMessage GlobalMessageStruct
		if err := json.Unmarshal(message, &pMessage); err != nil {
			log.Printf("Error parsing global message: %s | Data: %s\n", err, string(message)[:100])
			return err
		}
		tkrData, trdData, err := processMessageType(pMessage.Data.EventType, message, true)
		if err != nil {
			log.Printf("Error Processing Message Type %s from exchange %s.\nError Code: %s", pMessage.Data.EventType, exchange.Name, err)
			return err
		}
		if tkrData != (structs.TickerData{}) {
			// write something here
		} else if trdData != (structs.TradeData{}) {
			// write something here
		}
	}
	return nil
}

// Handle Ticker Messages
func HandleTickerMessage(message []byte, wrapped bool) (structs.TickerData, error) {
	var pData structs.TickerData
	if wrapped {
		// code this in
	} else {
		// code this in
	}
	return pData, nil
}

// Handle Trade Messages
func HandleTradeMessage(message []byte, wrapped bool) (structs.TradeData, error) {
	var pData structs.TradeData
	if wrapped {
		// code this in
	} else {
		// code this in
	}
	return pData, nil
}

// Gracefully close the connection
func CloseConnection(conn *websocket.Conn) {
	if conn == nil {
		log.Println("CloseConnection called with nil connection")
		return
	}
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Normal closure")
	if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
		log.Printf("Error sending close message: %v", err)
	}
	time.Sleep(1 * time.Second)
	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	} else {
		log.Println("Connection closed gracefully")
	}
}
