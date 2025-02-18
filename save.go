package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Antkky/go_crypto_scraper/structs"
)

type BufferedSaver struct {
	buffer      []structs.TickerData
	bufferSize  int
	mutex       sync.Mutex
	flushTicker *time.Ticker
}

func NewBufferedSaver(bufferSize int, flushInterval time.Duration) *BufferedSaver {
	saver := &BufferedSaver{
		buffer:      make([]structs.TickerData, 0, bufferSize),
		bufferSize:  bufferSize,
		flushTicker: time.NewTicker(flushInterval),
	}

	// Start background flushing process
	go saver.startFlusher()
	return saver
}

func (bs *BufferedSaver) AppendData(data structs.TickerData) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()

	bs.buffer = append(bs.buffer, data)
	if len(bs.buffer) >= bs.bufferSize {
		bs.flush()
	}
}

func (bs *BufferedSaver) flush() {
	if len(bs.buffer) == 0 {
		return
	}

	log.Println("Flushing buffer to file...")
	file, err := os.OpenFile("ticker_data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	bs.mutex.Lock()
	for _, data := range bs.buffer {
		if err := encoder.Encode(data); err != nil {
			log.Printf("Error writing to file: %s\n", err)
		}
	}
	bs.buffer = bs.buffer[:0] // Clear buffer after writing
	bs.mutex.Unlock()
}

func (bs *BufferedSaver) startFlusher() {
	for range bs.flushTicker.C {
		bs.flush()
	}
}

func (bs *BufferedSaver) Stop() {
	bs.flushTicker.Stop()
	bs.flush()
}
