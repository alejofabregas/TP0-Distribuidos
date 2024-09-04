package common

import (
	"fmt"
	"os"
)

// Bet structure used by the client
type Bet struct {
	AgencyID    string
	FirstName   string
	LastName    string
	DocumentID  string
	BirthDate   string
	BetNumber   string
}

// NewBetFromEnv initializes a new Bet structure from environment variables
func NewBetFromEnv() (*Bet, error) {
	bet := &Bet{
		AgencyID:    os.Getenv("CLI_ID"),
		FirstName:   os.Getenv("NOMBRE"),
		LastName:    os.Getenv("APELLIDO"),
		DocumentID:  os.Getenv("DOCUMENTO"),
		BirthDate:   os.Getenv("NACIMIENTO"),
		BetNumber:   os.Getenv("NUMERO"),
	}

	// Check that all environment variables are defined
	if bet.AgencyID == "" || bet.FirstName == "" || bet.LastName == "" ||
		bet.DocumentID == "" || bet.BirthDate == "" || bet.BetNumber == "" {
		return nil, fmt.Errorf("One or more environment variables are undefined")
	}

	return bet, nil
}

// ToBytes transforms a Bet to a byte slice
/*func (b *Bet) ToBytes() []byte {
	message := []byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s\n", b.AgencyID, b.FirstName, b.LastName, b.DocumentID, b.BirthDate, b.BetNumber))
	length := uint32(len(message))
	lengthBytes := make([]byte, 4) // 32 bits == 4 bytes
	binary.BigEndian.PutUint32(lengthBytes, length)
	result := append(lengthBytes, message...)
	return result
}*/

func (b *Bet) ToBytes() []byte {
	return []byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s\n", b.AgencyID, b.FirstName, b.LastName, b.DocumentID, b.BirthDate, b.BetNumber))
	/*length := uint32(len(message))
	lengthBytes := make([]byte, 4) // 32 bits == 4 bytes
	binary.BigEndian.PutUint32(lengthBytes, length)
	result := append(lengthBytes, message...)
	return result*/
}