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

func WrappedCheck(message []byte) (bool, error) {
	var pMessage GlobalMessageStruct

	if err := json.Unmarshal(message, &pMessage); err != nil {
		return false, err
	}

	if pMessage.Data.EventType != "" {
		return true, nil
	}
	if pMessage.EventType != "" {
		return false, nil
	}
	return false, errors.New("unknown message type")
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

func extractEventType(msg GlobalMessageStruct) string {
	if msg.Data.EventType != "" {
		return msg.Data.EventType
	}
	return msg.EventType
}

func processWrapped(wrapped bool, message []byte, bmessage *[]byte) error {
	if wrapped {
		var wrappedMsg struct {
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(message, &wrappedMsg); err != nil {
			return err
		}
		*bmessage = wrappedMsg.Data
	} else {
		*bmessage = message
	}
	return nil
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
func ProcessMessage(message []byte, tickerDataP *utils.TickerDataStruct, tradeData *utils.TradeDataStruct) (int, error) {
	var pMessage GlobalMessageStruct

	if err := json.Unmarshal(message, &pMessage); err != nil {
		return 0, err
	}

	switch extractEventType(pMessage) {
	case "24hrTicker":
		var bmessage []byte
		var tickerMsg TickerData
		wrapped, err := WrappedCheck(message)
		if err != nil {
			return 1, err
		}

		if err := processWrapped(wrapped, message, &bmessage); err != nil {
			return 1, err
		}

		if err := json.Unmarshal(bmessage, &tickerMsg); err != nil {
			return 1, err
		}

		*tickerDataP = utils.TickerDataStruct{
			TimeStamp: uint64(tickerMsg.EventTime),
			Date:      uint64(tickerMsg.EventTime),
			Symbol:    tickerMsg.Symbol,
			BidPrice:  string(tickerMsg.BidPrice),
			BidSize:   string(tickerMsg.BidSize),
			AskPrice:  string(tickerMsg.AskPrice),
			AskSize:   string(tickerMsg.AskSize),
		}

		return 1, nil
	case "trade":
		var bmessage []byte
		var tradeMsg TradeData
		wrapped, err := WrappedCheck(message)
		if err != nil {
			return 2, err
		}

		if err := processWrapped(wrapped, message, &bmessage); err != nil {
			return 2, err
		}

		if err := json.Unmarshal(bmessage, &tradeMsg); err != nil {
			return 2, err
		}

		*tradeData = utils.TradeDataStruct{
			TimeStamp: uint64(tradeMsg.EventTime),
			Date:      uint64(tradeMsg.EventTime),
			Symbol:    tradeMsg.Symbol,
			Price:     tradeMsg.Price,
			Quantity:  tradeMsg.Quantity,
			Bid_MM:    tradeMsg.IsMaker,
		}
		return 2, nil
	default:
		return 0, errors.New("unknown message type")
	}
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

	dataType, err := ProcessMessage(message, &tickerData, &tradeData)
	if err != nil {
		return err
	}
	if dataType == 0 {
		return errors.New("unknown message type")
	}

	if dataType == 1 && tickerData.Symbol != "" {
		log.Printf("Ticker data for %s", exchange.Name)
	} else if dataType == 2 && tradeData.Symbol != "" {
		log.Printf("Trade data for %s", exchange.Name)
	} else {
		return errors.New("data is possibly empty: " + exchange.Name)
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
