package common

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"os"
	"os/signal"
	"syscall"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

func (c *Client) StartClientLoop() {
	
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGINT)

	go func() {
		signal := <-signalChan
		
		log.Infof("action: client_shutdown | result: in_progress | client_id: %v | signal: %v", c.config.ID, signal)
		
		close(signalChan)
		if c.conn != nil {
			c.conn.Close()
		}
		
		log.Infof("action: client_shutdown | result: success | client_id: %v", c.config.ID)
		os.Exit(0)
	}()

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// Create Bet from environment variables
		bet, err := NewBetFromEnv()
		if err != nil {
			fmt.Println("Error creating Bet structure:", err)
			return
		}

		// TODO: Modify the send to avoid short-write
		// Send Bet to server through socket
		n, err := c.WriteAll(bet.ToBytes())
		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			bet.DocumentID,
			bet.BetNumber,
		)
		if n < len(bet.ToBytes()) {
			log.Errorf("action: apuesta_enviada | result: fail | short_read")
		}

		// Receive Bet ACK from server
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

		// Close the connection to the server
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			string(msg),
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) WriteAll(buffer []byte) (int, error) {
	bytesWritten := 0
	for bytesWritten < len(buffer) {
		// Escribir los bytes restantes del buffer
		n, err := c.conn.Write(buffer[bytesWritten:])
		if err != nil {
			return bytesWritten, err
		}
		bytesWritten += n
	}
	return bytesWritten, nil
}
