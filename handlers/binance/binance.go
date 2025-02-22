package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/gorilla/websocket"
)

// ________Small Helper Functions________

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

// ________Main Functions________

// InitializeStreams()
//
// Inputs:
//
//	conn        : *websocket.Conn
//	exchange    : utils.ExchangeConfig
//	dataBuffers : []*utils.DataBuffer
//
// Outputs:
//
//	No Outputs
//
// Description:
//
//	Initializes the streams by subscribing to the streams and creating the data buffers
func InitializeStreams(conn *websocket.Conn, exchange utils.ExchangeConfig, dataBuffers *map[string]*utils.DataBuffer) error {
	*dataBuffers = make(map[string]*utils.DataBuffer) // Initialize the map

	for _, stream := range exchange.Streams {
		bMessage, err := json.Marshal(stream.Message)
		filename := fmt.Sprintf("%s_%s_%s.csv", strings.ReplaceAll(exchange.Name, " ", ""), stream.Symbol, stream.Type)
		bufferCode := fmt.Sprintf("%s:%s@%s", stream.Symbol, stream.Type, strings.ReplaceAll(exchange.Name, " ", ""))
		filePath := fmt.Sprintf("data/%s/%s", strings.ReplaceAll(exchange.Name, " ", ""), stream.Symbol)
		(*dataBuffers)[bufferCode] = utils.NewDataBuffer(stream.Type, stream.Market, bufferCode, 250, filename, filePath)

		if err != nil {
			log.Printf("❌ Error marshalling subscribe message %v: %s", stream, err)
			return err
		}

		if err := conn.WriteMessage(websocket.TextMessage, bMessage); err != nil {
			log.Printf("❌ Error subscribing to stream %v: %s", stream, err)
			return err
		}
	}
	return nil
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
func HandleConnection(conn *websocket.Conn, exchange utils.ExchangeConfig) {
	if conn == nil {
		log.Println("Connection is nil, exiting HandleConnection.")
		return
	}

	// isEmpty checks if the given data is empty

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	messageQueue := make(chan []byte, 500)
	done := make(chan struct{})

	dataBuffers := make(map[string]*utils.DataBuffer)
	if err := InitializeStreams(conn, exchange, &dataBuffers); err != nil {
		return
	}

	go ConsumeMessages(messageQueue, exchange, done, dataBuffers)
	go ReceiveMessages(conn, messageQueue, done, exchange)

	<-interrupt
	log.Println("Interrupt received, closing connection...")

	CloseConnection(conn, exchange.Name)
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
	wrapped, err := WrappedCheck(message)
	if err != nil {
		return 0, err
	}

	var bmessage []byte
	if err := processWrapped(wrapped, message, &bmessage); err != nil {
		return 0, err
	}

	if err := json.Unmarshal(bmessage, &pMessage); err != nil {
		return 0, err
	}

	switch extractEventType(pMessage) {
	case "24hrTicker":
		var tickerMsg TickerData
		if err := json.Unmarshal(bmessage, &tickerMsg); err != nil {
			return 1, err
		}
		*tickerDataP = utils.TickerDataStruct{
			TimeStamp: uint64(tickerMsg.EventTime),
			Symbol:    tickerMsg.Symbol,
			BidPrice:  string(tickerMsg.BidPrice),
			BidSize:   string(tickerMsg.BidSize),
			AskPrice:  string(tickerMsg.AskPrice),
			AskSize:   string(tickerMsg.AskSize),
		}
		return 1, nil

	case "trade":
		var tradeMsg TradeData
		if err := json.Unmarshal(bmessage, &tradeMsg); err != nil {
			return 2, err
		}
		*tradeData = utils.TradeDataStruct{
			TimeStamp: uint64(tradeMsg.EventTime),
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
//	Processes incoming messages and adds them to the appropriate data buffer.
//	This function performs constant time lookups for the buffer associated with each message.
func ConsumeMessages(messageQueue chan []byte, exchange utils.ExchangeConfig, done chan struct{}, buffers map[string]*utils.DataBuffer) {
	defer close(done)
	normalizedExchangeName := strings.ReplaceAll(exchange.Name, " ", "")

	for message := range messageQueue {
		var (
			tickerData utils.TickerDataStruct
			tradeData  utils.TradeDataStruct
			bufferCode string
		)

		dataType, err := ProcessMessage(message, &tickerData, &tradeData)
		if err != nil {
			log.Printf("❌ Error processing message: %v", err)
			continue
		}

		switch dataType {
		case 0:
			log.Println("⚠️ Unknown message type, skipping message")
			continue
		case 1:
			if tickerData.Symbol != "" {
				bufferCode = fmt.Sprintf("%s:ticker@%s", tickerData.Symbol, normalizedExchangeName)
			}
		case 2:
			if tradeData.Symbol != "" {
				bufferCode = fmt.Sprintf("%s:trade@%s", tradeData.Symbol, normalizedExchangeName)
			}
		}

		if bufferCode != "" {
			if buffer, exists := buffers[bufferCode]; exists {
				if dataType == 1 {
					if err := buffer.AddData(tickerData); err != nil {
						log.Println("Error adding data to buffer: ", err)
						return
					}
				} else {
					if err := buffer.AddData(tradeData); err != nil {
						log.Println("Error adding data to buffer: ", err)
						return
					}
				}
			} else {
				log.Printf("⚠️ No buffer found for ID: %s", bufferCode)
			}
		}
	}
}

// ReceiveMessages()
//
// Inputs:
//
//	conn          : *websocket.Conn
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
//	Reads messages from the WebSocket connection and sends them to the messageQueue channel.
func ReceiveMessages(conn *websocket.Conn, messageQueue chan []byte, done chan struct{}, exchange utils.ExchangeConfig) {
	defer close(done)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from %s: %v", exchange.Name, err)
			return
		}

		select {
		case messageQueue <- message:
			//bruh
		case <-time.After(time.Millisecond * 100):
			log.Println("Producer slowed down")
		default:
			log.Printf("⚠️ Message queue full, dropping message for %s", exchange.Name)
		}
	}
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
