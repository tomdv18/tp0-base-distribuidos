package common

import (
	"bufio"
	"net"
	"io"
	"encoding/binary"
	"fmt"
	"strings"
)


func send_message(conn net.Conn, id string, bets []ClientData) (string, error) {
	var messages []string
	var msg string

	// Crear el array de mensajes
	for _, bet := range bets {
		message := fmt.Sprintf("%s;%s;%s;%s;%s;%s", id, bet.Nombre, bet.Apellido, bet.Documento, bet.Nacimiento, bet.Numero)
		messages = append(messages, message)
	}

	log.Infof("action: ammount_messages | result: success | messages: %v", len(messages))


		chunkMessage := strings.Join(messages, ">")

		// Verifico que el mensaje no exceda el tamaño permitido (8 kbytes)
		if len(chunkMessage) > 8192 {
			return "", fmt.Errorf("The message is too long")
		}

		// Envio el largo
		err := binary.Write(conn, binary.BigEndian, uint16(len(chunkMessage)))
		if err != nil {
			return "", fmt.Errorf("failed to write message length: %v", err)
		}
		//envio el mensaje
		_, err = io.WriteString(conn, chunkMessage)
		if err != nil {
			return "", fmt.Errorf("failed to send message: %v", err)
		}
		// Recibo respuesta del servidor
		msg, err = bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Errorf("action: receive_message | result: fail | error: %v", err)
			return "", err
		}

		log.Infof("action: receive_message | result: success ")
	
	return msg, nil
}
