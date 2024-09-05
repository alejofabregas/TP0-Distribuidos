package common

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// Betreader structure that stores the reading state of the CSV file.
type BetReader struct {
	agencyID string
	reader   *csv.Reader
	file     *os.File
	batchMaxAmount int
}

// NewBeatReader initializes a BetReader that opens and reads from the CSV file.
func NewBetReader(batchMaxAmount int) (*BetReader, error) {
	agencyID := os.Getenv("CLI_ID")
	if agencyID == "" {
		return nil, fmt.Errorf("Undefined environment variable CLI_ID")
	}

	filePath := fmt.Sprintf("./agency-%s.csv", agencyID)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error while opening file: %v", err)
	}

	reader := csv.NewReader(file)
	return &BetReader{
		agencyID: agencyID,
		reader:   reader,
		file:     file,
		batchMaxAmount: batchMaxAmount,
	}, nil
}

// Close frees the file resources after reading all batches.
func (br *BetReader) Close() {
	br.file.Close()
}

// NewBetBatch reeds a batch of 135 lines from the CSV file and returns a byte buffer
// and the total length of the buffer. Reads from last point read.
func (br *BetReader) NewBetBatch() ([]byte, uint32, bool, error) {
	var buffer bytes.Buffer
	lineCount := 0
	eofReached := false

	for lineCount < br.batchMaxAmount {
		record, err := br.reader.Read()
		if err == io.EOF {
			eofReached = true
			break
		}
		if err != nil {
			return nil, 0, eofReached, fmt.Errorf("error al leer el archivo CSV: %v", err)
		}

		if len(record) < 5 {
			return nil, 0, eofReached, fmt.Errorf("formato incorrecto en la línea %d", lineCount+1)
		}

		bet := Bet{
			AgencyID:   br.agencyID,
			FirstName:  record[0],
			LastName:   record[1],
			DocumentID: record[2],
			BirthDate:  record[3],
			BetNumber:  record[4],
		}

		buffer.Write(bet.ToBytes())
		lineCount++
	}

	return buffer.Bytes(), uint32(buffer.Len()), eofReached, nil
}
