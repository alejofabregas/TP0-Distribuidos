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
func (b *Bet) ToBytes() []byte {
	/*var buffer bytes.Buffer

	buffer.WriteString(b.AgencyID)
	buffer.WriteString("|") // Delimitador entre campos
	buffer.WriteString(b.FirstName)
	buffer.WriteString("|")
	buffer.WriteString(b.LastName)
	buffer.WriteString("|")
	buffer.WriteString(b.DocumentID)
	buffer.WriteString("|")
	buffer.WriteString(b.BirthDate)
	buffer.WriteString("|")
	buffer.WriteString(b.BetNumber)
	buffer.WriteString("\n")

	return buffer.Bytes()*/
	//var bytesBuffer []byte
	string := fmt.Sprintf("%s|%s|%s|%s|%s|%s\n", b.AgencyID, b.FirstName, b.LastName, b.DocumentID, b.BirthDate, b.BetNumber)
	return []byte(string)
}

/*func main() {
	// Ejemplo de uso
	bet, err := NewBetFromEnv()
	if err != nil {
		fmt.Println("Error creating Bet structure:", err)
		return
	}

	bytes := bet.ToBytes()

	fmt.Println("Bet bytes:", bytes)
	fmt.Println("Bet as string:", string(bytes))
}*/