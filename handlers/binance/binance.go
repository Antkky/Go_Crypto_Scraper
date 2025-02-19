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

// HandleConnection()
//
// Inputs:
//
//	conn     : *websocket.Conn
//	exchange : utils.ExchangeConfig
//
// Outputs:
//
//	No Outputs
//
// Description:
//
//	goroutine that subscribes and launches 2 goroutines to listen for messages and handle them
func HandleConnection(conn *websocket.Conn, exchange utils.ExchangeConfig, buffer utils.DataBuffer) {
	if conn == nil {
		log.Println("Connection is nil, exiting HandleConnection.")
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	messageQueue := make(chan []byte, 100)
	done := make(chan struct{})

	go ConsumeMessages(messageQueue, done, exchange)
	go ReceiveMessages(conn, messageQueue, done, exchange)

	<-interrupt
	log.Println("Interrupt received, closing connection...")

	CloseConnection(conn, exchange.Name)
}

// ConsumeMessages()
//
// Inputs:
//
//	messageQueue  : chan []byte
//	done         : chan struct{}
//	exchange     : utils.ExchangeConfig
//
// Outputs:
//
//	No Outputs
//
// Description:
//
//	invoke the HandleMessage() function to process the message in the messageQueue
func ConsumeMessages(messageQueue chan []byte, done chan struct{}, exchange utils.ExchangeConfig) {
	for message := range messageQueue {
		if err := HandleMessage(message, exchange); err != nil {
			log.Printf("Error handling message for %s: %v", exchange.Name, err)
		}
	}
	close(done)
}

// ReceiveMessages()
//
// Inputs:
//
//	message      : *websocket.Conn
//	messageQueue : chan []byte,
//	done         : chan struct{},
//	exchange     : utils.ExchangeConfig
//
// Outputs:
//
//	No Outputs
//
// Description:
//
//	sends received messages to the messageQueue channel
func ReceiveMessages(conn *websocket.Conn, messageQueue chan []byte, done chan struct{}, exchange utils.ExchangeConfig) {
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

// ProcessMessageType()
//
// Inputs:
//
//	message    : []byte
//	tickerData : *utils.TickerDataStruct
//	tradeData  : *utils.TradeDataStruct
//
// Outputs:
//
//	error
//
// Description:
//
//	basically routes the data to the correct processing function
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
		if err := json.Unmarshal(message, tickerData); err != nil {
			return err
		}
	case "trade":
		if err := json.Unmarshal(message, tradeData); err != nil {
			return err
		}
	default:
		return errors.New("unhandled event type: " + eventType)
	}
	return nil
}

// HandleMessage()
//
// Inputs:
//
//	message  : []byte
//	exchange : utils.ExchangeConfig
//
// Outputs:
//
//	error
//
// Description:
//
//	handle processing and saving the message
func HandleMessage(message []byte, exchange utils.ExchangeConfig) error {
	var (
		tickerData utils.TickerDataStruct
		tradeData  utils.TradeDataStruct
	)

	if err := ProcessMessageType(message, &tickerData, &tradeData); err != nil {
		return err
	}

	if (tickerData != utils.TickerDataStruct{}) {
		log.Printf("Ticker data for %s", exchange.Name)
	} else if (tradeData != utils.TradeDataStruct{}) {
		log.Printf("Trade data for %s", exchange.Name)
	}

	return nil
}

// Subscribe()
//
// Inputs:
//
//	conn         : *websocket.Conn
//	stream       : map[string]interface{}
//	exchangeName : string
//
// Outputs:
//
//	error
//
// Description:
//
//	sends a subscribe message to the server
func Subscribe(conn *websocket.Conn, stream map[string]interface{}) error {
	message, err := json.Marshal(stream)
	if err != nil {
		return err
	}
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}
	return nil
}

// CloseConnection()
//
// Inputs:
//
//	conn         : *websocket.conn
//	exchangeName : string
//
// Outputs:
//
//	No Outputs
//
// Description:
//
//	Gracefully close the connection by sending a closure message and gracefully close connection
func CloseConnection(conn *websocket.Conn, exchangeName string) {
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Normal closure")
	if err := conn.WriteMessage(websocket.CloseMessage, closeMsg); err != nil {
		log.Printf("Error sending close message for %s: %v", exchangeName, err)
	}

	time.Sleep(time.Second)

	if err := conn.Close(); err != nil {
		log.Printf("Error closing connection for %s: %v", exchangeName, err)
	} else {
		log.Printf("Connection for %s closed gracefully", exchangeName)
	}
}
