package common

import (
	"bufio"
	"net"
	"io"
	"encoding/binary"
	"fmt"
	"errors"
)
func send_message(conn net.Conn, id string, bets []ClientData, maxBatchSize int) (string, error) {
	var messages []string
	var err error
	var response string

	for _, bet := range bets {
		message := fmt.Sprintf("%s;%s;%s;%s;%s;%s", id, bet.Nombre, bet.Apellido, bet.Documento, bet.Nacimiento, bet.Numero)
		messages = append(messages, message)
	}

	log.Infof("action: ammount_messages | result: success | messages_count: %d", len(messages))

	for len(messages) > 0 {
		var batch []string
		var batchMessage string

		// Tomar hasta maxBatchSize elementos o menos si quedan pocos
		batchSize := maxBatchSize
		if len(messages) < maxBatchSize {
			batchSize = len(messages)
		}
		batch = messages[:batchSize]

		batchMessage = strings.Join(batch, ">")
		messages = messages[batchSize:]

		// Verificar que no supere 8KB
		if len(batchMessage) > 8192 {
			return "", errors.New("The message size is too large")
		}
		binary.Write(conn, binary.BigEndian, uint16(len(batchMessage)))
		_, err = io.WriteString(conn, batchMessage)
		if err != nil {
			log.Errorf("action: send_message | result: fail | error: %v", err)
			return "", err
		}

		log.Infof("action: send_message | result: success | batch_size: %d | message: %v", batchSize, batchMessage)


	}

	response, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Errorf("action: receive_message | result: fail | error: %v", err)
		return "", err
	}

	log.Infof("action: receive_message | result: success | server_response: %v", response)

	return response, err
}
