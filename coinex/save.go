package coinex

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var batchSize uint = 256

type dataStore struct {
	sync.RWMutex
	data map[string][][]string
}

var (
	coinexTickerData = &dataStore{data: make(map[string][][]string)}
	coinexTradeData = &dataStore{data: make(map[string][][]string)}
	coinexLastFlushTime sync.Map
)

func updateStoredRows(symbol string, store *dataStore, rows ...[]string) {
	store.Lock()
	defer store.Unlock()

	if store.data[symbol] == nil {
			store.data[symbol] = make([][]string, 0)
	}
	store.data[symbol] = append(store.data[symbol], rows...)
}

func flushRowsToCSV(symbol string, exchange string, filename string, store *dataStore, filetype string) error {
	dir := fmt.Sprintf("data/%s/%s", exchange, symbol)
	filePath := fmt.Sprintf("%s/%s", dir, filename)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Println("Error creating directory:", err)
		return err
	}

	rows := store.data[symbol]
	rowCount := len(rows)

	if rowCount == 0 {
		return nil
	}

	lastFlushTimeVal, _ := coinexLastFlushTime.LoadOrStore(symbol, time.Time{})
  lastFlushTime, _ := lastFlushTimeVal.(time.Time)

	if rowCount >= int(batchSize) || time.Since(lastFlushTime) > 5*time.Second {
		if err := writeRowsToCSV(filePath, rows, filetype); err != nil {
				log.Println("Error writing CSV:", err)
				return err
		}

		store.Lock()
		store.data[symbol] = make([][]string, 0)
		store.Unlock()

		coinexLastFlushTime.Store(symbol, time.Now())
	}

	return nil
}

func writeRowsToCSV(filePath string, rows [][]string, fileType string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening CSV file:", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var header []string
	if fileType == "Ticker" {
		header = []string{"TimeStamp", "Date"}//Header here
	} else if fileType == "Trades" {
		header = []string{"TimeStamp", "Date"}//Header here
	}

	fileStats, err := file.Stat()
	if err != nil {
		log.Panicln("Error getting file stats: ", err)
	}
	if fileStats.Size() == 0 {
		if err := writer.Write(header); err != nil {
			log.Println("Error writing header to CSV: ", err)
			return err
		}
	}

	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			log.Println("Error writing row to CSV file: ", err)
			return err
		}
	}

	return nil
}
