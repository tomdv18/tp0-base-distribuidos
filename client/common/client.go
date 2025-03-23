package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	config  ClientConfig
	conn    net.Conn
	sigChan chan os.Signal
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, sigChan chan os.Signal) *Client {
	return &Client{
		config:  config,
		sigChan: sigChan,
	}
}

// createClientSocket Initializes client socket.
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}

// closeClientSocket Closes the client socket gracefully
func (c *Client) closeClientSocket() {
	if c.conn != nil {
		log.Infof("Closing client socket for client_id: %v", c.config.ID)
		c.conn.Close()
	}
}

// StartClientLoop Handles client message sending with graceful shutdown
func (c *Client) StartClientLoop() {
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Manejo de interrupción SIGTERM
		select {
		case <-c.sigChan:
			log.Infof("Received SIGTERM. Shutting down gracefully...")
			c.closeClientSocket()
			os.Exit(0)
		default:
		}

		// Crear conexión al servidor
		err := c.createClientSocket()
		if err != nil {
			return
		}

		// Enviar mensaje
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)

		// Leer respuesta del servidor
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.closeClientSocket() // Cerrar conexión después de cada iteración

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Esperar antes de enviar el siguiente mensaje
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
