package common
import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
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
	config ClientConfig
	conn   net.Conn
	quit chan struct {}
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		quit: make(chan struct{}),
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	go c.handleShutdown() // Capturar SIGTERM o SIGINT en otra goroutine

	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Si se recibió la señal de salida, terminamos el bucle
		select {
		case <-c.quit:
			log.Infof("action: client_exit | result: received_shutdown_signal | client_id: %v", c.config.ID)
			return
		default:
		}

		err := c.createClientSocket()
		if err != nil {
			log.Errorf("action: connect | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		fmt.Fprintf(c.conn, "[CLIENT %v] Message N°%v\n", c.config.ID, msgID)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v", c.config.ID, msg)

		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}



func (c *Client) handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	log.Infof("action: shutdown | result: received_signal | client_id: %v", c.config.ID)

	select {
	case <-c.quit:
		
	default:
		close(c.quit) // Solo se cierrar si todavia no se cerro
	}
	c.closeClientSocket() 

	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)

	os.exit(0)
}

func (c *Client) closeClientSocket() {
	if c.conn != nil {
		log.Infof("action: close_socket | result: success | client_id: %v", c.config.ID)
		_ = c.conn.Close() // Ignorar el error si la conexión ya está cerrada
		c.conn = nil
	}
}



