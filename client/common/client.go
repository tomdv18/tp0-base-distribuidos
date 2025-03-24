package common

import (
	"net"
	"time"
	"os"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BachMaxAmmount  int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	quit chan os.Signal
	clientData []ClientData
}

type ClientData struct {
	Nombre string
	Apellido string
	Documento string
	Nacimiento string
	Numero string
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, quit chan os.Signal, clientData []ClientData) *Client {
	client := &Client{
		config: config,
		quit: quit,
		clientData: clientData,
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


func (c *Client) shutdown_client() {
	c.conn.Close()
	log.Infof("action: socket_closing | result: success | client_id: %v",c.config.ID)
}


// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {

	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()


		msg, err := send_message(c.conn, c.config.ID, c.clientData, c.config.BachMaxAmmount)
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
			c.clientData[0].Documento,
			msg,
		)


		select	{
		case <-c.quit:
			log.Infof("action: finish_signal | result: in_progress | client_id: %v", c.config.ID)
			c.shutdown_client()
			log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
			return
		default:
		}
		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	time.Sleep(0100 * time.Millisecond)
}