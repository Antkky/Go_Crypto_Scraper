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



//___________Table Of Contents______________
// func HandleConnection : line - 32
// func ParseStreams : line - 78
// func Subscribe : line - 87
// func ConsumeMessages : line - 124
// func ReceiveMessages : line - 147
// func ProcessMessageType : line - 179
// func HandleTickerMessage : line - 212
// func HandleTradeMessage : line - 231
// func HandleMessage : line - 249
// func CloseConnection : line - 276




// HandleConnection()
//
// Inputs:
//  conn     : *websocket.Conn
//  exchange : utils.ExchangeConfig
//
// Outputs:
//  No Outputs
//
// Description:
//  goroutine that subscribes and launches 2 goroutines to listen for messages and handle them
//
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

  for _, stream := range streams {
    if err := Subscribe(conn, stream, exchange.Name); err != nil {
      log.Printf("Error subscribing to stream for %s: %s", exchange.name, err)
    }
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

	closeConnection(conn, exchange.Name)
}

// ParseStreams()
//
// Inputs:
//  streamsData : json.RawMessage
//
// Outputs:
//  []map[string]interface{}
//  error
//
// Description:
//  turns the streamsData into a data type we can iterate over
//
func ParseStreams(streamsData json.RawMessage) ([]map[string]interface{}, error) {
	var streams []map[string]interface{}
	if err := json.Unmarshal(streamsData, &streams); err != nil {
		return nil, err
	}
	return streams, nil
}


// Subscribe()
//
// Inputs:
//  conn         : *websocket.Conn
//  stream       : map[string]interface{}
//  exchangeName : string
//
// Outputs:
//  error
//
// Description:
//  sends a subscribe message to the server
//
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


// ConsumeMessages()
//
// Inputs:
//  messageQueue  : chan []byte
//  done         : chan struct{}
//  exchange     : utils.ExchangeConfig
//
// Outputs:
//  No Outputs
//
// Description:
//  invoke the HandleMessage() function to process the message in the messageQueue
//
func ConsumeMessages(messageQueue chan []byte, done chan struct{}, exchange utils.ExchangeConfig) {
	defer close(done)
	for message := range messageQueue {
		if err := HandleMessage(message, exchange); err != nil {
			log.Printf("Error handling message for %s: %v", exchange.Name, err)
		}
	}
}


// ReceiveMessages()
//
// Inputs:
//  message      : *websocket.Conn
//  messageQueue : chan []byte,
//  done         : chan struct{},
//  exchange     : utils.ExchangeConfig
//
// Outputs:
//  No Outputs
//
// Description:
//  sends received messages to the messageQueue channel
//
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


// ProcessMessageType()
//
// Inputs:
//  message    : []byte
//  tickerData : *utils.TickerDataStruct
//  tradeData  : *utils.TradeDataStruct
//
// Outputs:
//  error
//
// Description:
//  basically routes the data to the correct processing function
//
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

// HandleTickerMessage()
//
// Inputs:
//  message   : []byte
//  tradeData : *utils.TradeDataStruct
//
// Outputs:
//  error
//
// Description:
//  processes the ticker data inside the byte array and push the data to the pointer tradeData
//
func HandleTickerMessage(message []byte, tickerData *utils.TickerDataStruct) error {
	// Add logic here to handle ticker message
	// e.g., Unmarshal and process the data, based on whether it's wrapped or not
	return nil
}


// HandleTradeMessage()
//
// Inputs:
//  message   : []byte
//  tradeData : *utils.TradeDataStruct
//
// Outputs:
//  error
//
// Description:
//  processes the trade data inside the byte array
//
func HandleTradeMessage(message []byte, tradeData *utils.TradeDataStruct) error {
	// Add logic here to handle trade message
	// e.g., Unmarshal and process the data, based on whether it's wrapped or not
	return nil
}

// HandleMessage()
//
// Inputs:
//  message  : []byte
//  exchange : utils.ExchangeConfig
//
// Outputs:
//  error
//
// Description:
//  handle processing and saving the message
//
func HandleMessage(message []byte, exchange utils.ExchangeConfig) error {
	var (
		tickerData utils.TickerDataStruct
		tradeData  utils.TradeDataStruct
	)

	if err := ProcessMessageType(message, &tickerData, &tradeData); err != nil {
		return err
	}

  // implement saving logic

	return nil
}

// CloseConnection()
//
// Inputs:
//  conn         : *websocket.conn
//  exchangeName : string
//
// Outputs:
//  No Outputs
//
// Description:
//  Gracefully close the connection by sending a closure message and gracefully close connection
//
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
