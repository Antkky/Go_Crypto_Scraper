package binance

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/Antkky/go_crypto_scraper/structs"
	"github.com/gorilla/websocket"
)

/*
Binance Global and Binance US has different response formats

Subscribe to streams
Make Ping Goroutines
Start Receiving Messages
HandleMessages()
*/
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

	done := make(chan struct{})

	// Start receiving incoming messages
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error receiving message from exchange %s: %s\n", exchange.Name, err)
				return
			}
			go HandleMessage(message, exchange)
		}
	}()
	<-interrupt
	log.Println("Interrupt received, closing connection...")
	// Attempt graceful WebSocket close
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
		log.Printf("Error closing WebSocket connection for %s: %s\n", exchange.Name, err)
	}
	conn.Close()
}

func HandleMessage(message []byte, exchange structs.ExchangeConfig) {
	// does some processing and saves it to a CSV file
	// will write later
}

func CloseConnection(conn *websocket.Conn) {
	// gracefully close the connection with the exchange
}
