package coinex

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Antkky/go_crypto_scraper/utils"
	"github.com/Antkky/go_crypto_scraper/utils/buffer"
	"github.com/gorilla/websocket"
)

func decompressGzip(data []byte) ([]byte, error) {
	// Check for gzip magic numbers (0x1f 0x8b)
	if len(data) < 2 || data[0] != 0x1f || data[1] != 0x8b {
		return nil, fmt.Errorf("invalid gzip header")
	}

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer reader.Close()

	var decompressed bytes.Buffer
	if _, err := io.Copy(&decompressed, reader); err != nil {
		return nil, fmt.Errorf("failed to decompress data: %w", err)
	}

	return decompressed.Bytes(), nil
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
func InitializeStreams(conn *websocket.Conn, exchange utils.ExchangeConfig, dataBuffers *map[string]*buffer.DataBuffer, logger *log.Logger) error {
	*dataBuffers = make(map[string]*buffer.DataBuffer)

	for _, stream := range exchange.Streams {
		bMessage, err := json.Marshal(stream.Message)
		filename := fmt.Sprintf("%s_%s_%s.csv", strings.ReplaceAll(exchange.Name, " ", ""), stream.Symbol, stream.Type)
		bufferCode := fmt.Sprintf("%s:%s@%s", stream.Symbol, stream.Type, strings.ReplaceAll(exchange.Name, " ", ""))
		filePath := fmt.Sprintf("data/%s/%s", strings.ReplaceAll(exchange.Name, " ", ""), stream.Symbol)
		(*dataBuffers)[bufferCode] = buffer.NewDataBuffer(stream.Type, stream.Market, bufferCode, 50, filename, filePath)

		if err != nil {
			logger.Printf("❌ Error marshalling subscribe message %v: %s", stream, err)
			return err
		}

		if err := conn.WriteMessage(websocket.TextMessage, bMessage); err != nil {
			logger.Printf("❌ Error subscribing to stream %v: %s", stream, err)
			return err
		}
		time.Sleep(500 * time.Millisecond)
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
func HandleConnection(conn *websocket.Conn, exchange utils.ExchangeConfig, logger *log.Logger) {
	if conn == nil {
		logger.Println("Connection is nil, exiting HandleConnection.")
		return
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	messageQueue := make(chan []byte, 500)
	done := make(chan struct{})

	dataBuffers := make(map[string]*buffer.DataBuffer)
	if err := InitializeStreams(conn, exchange, &dataBuffers, logger); err != nil {
		return
	}

	go ConsumeMessages(messageQueue, exchange, done, dataBuffers, logger)
	go ReceiveMessages(conn, messageQueue, exchange, done, logger)

	<-interrupt
	logger.Println("Interrupt received, closing connection...")

	CloseConnection(conn, exchange.Name, logger)
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
//	For more details, see the [Obsidian Documentation](obsidian://open?vault=Go_crypto_scraper&file=handlers/coinex/ProcessMessage.md).
func ProcessMessage(message []byte, tickerDataP *[]utils.TickerDataStruct, tradeDataP *[]utils.TradeDataStruct) (int, error) {
	decompressed, err := decompressGzip(message)
	if err != nil {
		return 0, err
	}

	var pMessage GlobalMessageStruct
	if err := json.Unmarshal(decompressed, &pMessage); err != nil {
		return 0, fmt.Errorf("failed to unmarshal coinex message: %w", err)
	}

	switch {
	case pMessage.Method == "bbo.update":
		var tickerMsg TickerData
		if err := json.Unmarshal(decompressed, &tickerMsg); err != nil {
			return 1, fmt.Errorf("failed to unmarshal ticker data: %w", err)
		}
		*tickerDataP = append(*tickerDataP, utils.TickerDataStruct{
			TimeStamp: uint64(tickerMsg.Data.Updated_at),
			Symbol:    tickerMsg.Data.Market,
			BidPrice:  string(tickerMsg.Data.BidPrice),
			BidSize:   string(tickerMsg.Data.BidSize),
			AskPrice:  string(tickerMsg.Data.AskPrice),
			AskSize:   string(tickerMsg.Data.AskSize),
		})
		return 1, nil

	case pMessage.Method == "deals.update":
		var tradeMsg TradeData
		if err := json.Unmarshal(decompressed, &tradeMsg); err != nil {
			return 1, fmt.Errorf("failed to unmarshal trade data: %w", err)
		}
		for _, trade := range tradeMsg.Data.Deals {
			*tradeDataP = append(*tradeDataP, utils.TradeDataStruct{
				TimeStamp: uint64(trade.Created_at),
				Symbol:    tradeMsg.Data.Market,
				Price:     trade.Price,
				Quantity:  trade.Amount,
				Bid_MM:    trade.Side == "sell",
			})
		}
		return 2, nil

	case bytes.Equal(decompressed, []byte(`{"id":1,"code":0,"message":"OK"}`)):
		return 5, nil

	default:
		return 0, fmt.Errorf("unknown message type: %s", pMessage.Method)
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
func ConsumeMessages(messageQueue chan []byte, exchange utils.ExchangeConfig, done chan struct{}, buffers map[string]*buffer.DataBuffer, logger *log.Logger) {
	defer close(done)
	normalizedExchangeName := strings.ReplaceAll(exchange.Name, " ", "")

	for message := range messageQueue {
		var bufferCode string
		tickerData := make([]utils.TickerDataStruct, 0)
		tradeData := make([]utils.TradeDataStruct, 0)

		dataType, err := ProcessMessage(message, &tickerData, &tradeData)
		if err != nil {
			logger.Printf("❌ Error processing message: %v", err)
			continue
		}

		switch dataType {
		case 0:
			continue
		case 1:
			if len(tickerData) > 0 && tickerData[0].Symbol != "" {
				bufferCode = fmt.Sprintf("%s:ticker@%s", tickerData[0].Symbol, normalizedExchangeName)
			}
		case 2:
			if len(tradeData) > 0 && tradeData[0].Symbol != "" {
				bufferCode = fmt.Sprintf("%s:trade@%s", tradeData[0].Symbol, normalizedExchangeName)
			}
		case 5:
			logger.Println("✅ Subscribe Success")
			continue
		}
		if bufferCode != "" {
			if buffer, exists := buffers[bufferCode]; exists {
				if dataType == 1 {
					if err := buffer.AddData(tickerData); err != nil {
						logger.Println("❌ Error adding data to buffer: ", err)
						return
					}
				} else {
					if err := buffer.AddData(tradeData); err != nil {
						logger.Println("❌ Error adding data to buffer: ", err)
						return
					}
				}
			} else {
				logger.Printf("❌ No buffer found for ID: %s", bufferCode)
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
func ReceiveMessages(conn *websocket.Conn, messageQueue chan []byte, exchange utils.ExchangeConfig, done chan struct{}, logger *log.Logger) {
	defer close(done)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logger.Printf("Error reading message from %s: %v", exchange.Name, err)
			return
		}

		select {
		case messageQueue <- message:
			//bruh
		case <-time.After(time.Millisecond * 100):
			logger.Println("Producer slowed down")
		default:
			logger.Printf("⚠️ Message queue full, dropping message for %s", exchange.Name)
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
func CloseConnection(conn *websocket.Conn, exchangeName string, logger *log.Logger) {
	if err := conn.Close(); err != nil {
		logger.Printf("Error closing connection for %s: %v", exchangeName, err)
	} else {
		logger.Printf("Connection for %s closed gracefully", exchangeName)
	}
}
