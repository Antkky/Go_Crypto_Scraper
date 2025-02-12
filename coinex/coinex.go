package coinex

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
)

func isGzipped(message []byte) bool {
	return len(message) > 2 && message[0] == 0x1f && message[1] == 0x8b
}

func decompressGzip(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip decompression error: %w", err)
	}
	defer gzipReader.Close()

	var decompressedData bytes.Buffer
	if _, err := io.Copy(&decompressedData, gzipReader); err != nil {
		return nil, fmt.Errorf("gzip copy error: %w", err)
	}
	return decompressedData.Bytes(), nil
}

func HandleMessage(message []byte) error {
	if isGzipped(message) {
		// message is gzipped
		//decompressedMessage, err := decompressGzip(message)
		//if err != nil {
		//return err
		//}
	} else {
		// message isnt gzipped
	}

	return nil
}
