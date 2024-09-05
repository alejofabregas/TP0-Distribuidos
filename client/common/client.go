package common

import (
	"bufio"
	"encoding/binary"
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

	betReader, err := NewBetReader()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	defer betReader.Close()

	go func() {
		signal := <-signalChan
		
		log.Infof("action: client_shutdown | result: in_progress | client_id: %v | signal: %v", c.config.ID, signal)
		
		close(signalChan)
		if c.conn != nil {
			c.conn.Close()
		}
		
		betReader.Close()
		
		log.Infof("action: client_shutdown | result: success | client_id: %v", c.config.ID)
		os.Exit(0)
	}()

	for {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()
		data, totalBytes, eofReached, err := betReader.NewBetBatch()
		if err != nil {
			log.Errorf("action: batch_leido | result: fail")
			return
		}
		log.Infof("action: batch_leido | result: success")

		lengthBytes := make([]byte, 4) // 32 bits == 4 bytes
		binary.BigEndian.PutUint32(lengthBytes, totalBytes)
		batch := append(lengthBytes, data...)

		// Send Batch to server through socket
		n, err := c.WriteAll(batch)
		log.Infof("action: batch_enviado | result: success")
		if n < len(batch) {
			log.Errorf("action: batch_enviado | result: fail | short_write")
		}

		// Receive Batch ACK from server
		msg, err := bufio.NewReader(c.conn).ReadString('\n')

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

		if eofReached {
			log.Infof("action: envio_completado | result: success")
			// Close the connection to the server
			c.conn.Close()
			break
		}
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
