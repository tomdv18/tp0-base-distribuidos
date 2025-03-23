package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
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

// Client Entity that encapsulates how the client works
type Client struct {
	config ClientConfig
	conn   net.Conn
	quit   chan os.Signal
}

// NewClient Initializes a new client receiving the configuration and quit channel
func NewClient(config ClientConfig, quit chan os.Signal) *Client {
	return &Client{
		config: config,
		quit:   quit,
	}
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}
	c.conn = conn
	log.Debugf("action: connect | result: success | client_id: %v", c.config.ID)
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	log.Infof("action: start_loop | result: success | client_id: %v", c.config.ID)

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Si recibimos una señal de terminación, salimos limpiamente
		select {
		case <-c.quit:
			log.Infof("action: received termination signal | result: shutting_down | client_id: %v", c.config.ID)
			if c.conn != nil {
				c.conn.Close()
			}
			return
		default:
			// Continuamos con el flujo normal
		}

		// Intentamos crear la conexión
		if err := c.createClientSocket(); err != nil {
			log.Errorf("action: retry | result: waiting | client_id: %v", c.config.ID)
			time.Sleep(2 * time.Second) // Esperamos antes de intentar nuevamente
			continue
		}

		// Enviar mensaje al servidor
		fmt.Fprintf(c.conn, "[CLIENT %v] Message N°%v\n", c.config.ID, msgID)

		// Leer respuesta
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v", c.config.ID, msg)

		// Esperar antes del siguiente mensaje
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
