package binance

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/gorilla/websocket"
)

func isEmpty(data interface{}) bool {
	switch v := data.(type) {
	case utils.TickerDataStruct:
		return v == utils.TickerDataStruct{}
	case utils.TradeDataStruct:
		return v == utils.TradeDataStruct{}
	default:
		return true
	}
}

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

	// isEmpty checks if the given data is empty

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
//	done          : chan struct{}
//	exchange      : utils.ExchangeConfig
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

func convertStringToFloat(value string) float32 {
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0 // Handle error properly in production
	}
	return float32(f)
}

func extractEventType(msg GlobalMessageStruct) string {
	if msg.Data.EventType != "" {
		return msg.Data.EventType
	}
	return msg.EventType
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
func ProcessMessage(message []byte, tickerData *utils.TickerDataStruct, tradeData *utils.TradeDataStruct) error {
	var pMessage GlobalMessageStruct

	if err := json.Unmarshal(message, &pMessage); err != nil {
		return err
	}


	switch extractEventType(pMessage) {
	case "24hrTicker":
		var tickerMsg TickerData
		if err := json.Unmarshal(message, &tickerMsg); err != nil {
			return err
		}
		*tickerData = utils.TickerDataStruct{
			TimeStamp: uint(tickerMsg.EventTime),
			Date:      uint(tickerMsg.EventTime),
			Symbol:    tickerMsg.Symbol,
			BidPrice:  tickerMsg.BidPrice.String(),
			BidSize:   tickerMsg.BidSize.String(),
			AskPrice:  tickerMsg.AskPrice.String(),
			AskSize:   tickerMsg.AskSize.String(),
		}
	case "trade":
		var tradeMsg utils.TradeDataStruct
		if err := json.Unmarshal(message, &tradeMsg); err != nil {
			return err
		}
		*tradeData = tradeMsg
	default:
		return errors.New("unknown message type")
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
	if err := ProcessMessage(message, &tickerData, &tradeData); err != nil {
		return err
	}

	// Handle ticker data
	if !isEmpty(tickerData) {
		log.Printf("Ticker data for %s: %+v", exchange.Name, tickerData)
	} else if !isEmpty(tradeData) {
		// Handle trade data
		log.Printf("Trade data for %s: %+v", exchange.Name, tradeData)
	} else {
		// Log when no useful data is found
		log.Printf("Received message for %s but data is empty: %s", exchange.Name, string(message))
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
