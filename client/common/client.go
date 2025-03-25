package common

import (
	"net"
	"time"
	"os"
	"github.com/op/go-logging"
	"strings"
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

func (c *Client) obtain_winners() {

	loop :
	for {
		c.createClientSocket()
		msg, err := send_winners(c.conn, c.config.ID)
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",	c.config.ID, err)
			c.conn.Close()
			return
		}
		msg = strings.TrimSpace(msg)

		if msg == "NOT_READY" {
			log.Infof("action: winners_received | result: in_progress | client_id: %v", c.config.ID)
			c.conn.Close()
			time.Sleep(c.config.LoopPeriod)
		} else {
			log.Infof("action: winners_received | result: success | winners: %v", msg)
			c.conn.Close()
			break loop
		}
		select	{
		case <-c.quit:
			log.Infof("action: finish_signal | result: in_progress | client_id: %v", c.config.ID)
			c.shutdown_client()
			break loop
		default:
		}
	}

}



// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	index := 0
    loop:
    for {
        if index >= len(c.clientData) {
            log.Infof("action: loop_finished | result: in_progress | client_id: %v", c.config.ID)
            break loop
        }

        // Crear la conexión al servidor en cada iteración del loop
        c.createClientSocket()

        end := index + c.config.BachMaxAmmount
        if end > len(c.clientData) {
            end = len(c.clientData)
        }

        chunk := c.clientData[index:end]

        // Enviar el chunk al servidor
        _, err := send_message(c.conn, c.config.ID, chunk)
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: apuesta_enviada | result: success")
		index = end

        c.conn.Close()


		select	{
		case <-c.quit:
			log.Infof("action: finish_signal | result: in_progress | client_id: %v", c.config.ID)
			c.shutdown_client()
			break loop
		default:
		}
		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	c.obtain_winners()
	time.Sleep(0100 * time.Millisecond)
}