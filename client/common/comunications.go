package common

import (
	"bufio"
	"net"
	"io"
	"encoding/binary"
	"fmt"
)


func send_message(conn net.Conn, id string, bets []ClientData, maxBatchSize int) (string, error) {
	var messages []string

	// Crear el array de mensajes
	for _, bet := range bets {
		message := fmt.Sprintf("%s;%s;%s;%s;%s;%s", id, bet.Nombre, bet.Apellido, bet.Documento, bet.Nacimiento, bet.Numero)
		messages = append(messages, message)
		log.Infof("action: create_message | result: success | message: %v", message)
	}

	log.Infof("action: ammount_messages | result: success | messages: %v", len(messages))

	
	for i := 0; i < len(messages); i += maxBatchSize {
		// Determinar el tamaño del batch (puede ser menor que maxBatchSize en el último grupo)
		end := i + maxBatchSize
		if end > len(messages) {
			end = len(messages)
		}

		// Crear un slice de los mensajes para este grupo
		batch := messages[i:end]
		
		// Concatenar los mensajes en el batch con ">" como delimitador, sin agregar ">" al final
		batchMessage := ""
		for j, msg := range batch {
			if j > 0 {
				batchMessage += ">"
			}
			batchMessage += msg
		}

		// Verificar que el mensaje no exceda el tamaño permitido (por ejemplo, 8192 bytes)
		if len(batchMessage) > 8192 {
			return "", fmt.Errorf("The message is too long")
		}

		// Enviar el mensaje
		err := binary.Write(conn, binary.BigEndian, uint16(len(batchMessage)))
		if err != nil {
			return "", fmt.Errorf("failed to write message length: %v", err)
		}
		_, err = io.WriteString(conn, batchMessage)
		if err != nil {
			return "", fmt.Errorf("failed to send message: %v", err)
		}

	}
	

	// Recibir respuesta del servidor
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Errorf("action: receive_message | result: fail | error: %v", err)
		return "", err
	}

	log.Infof("action: receive_message | result: success | server_response: %v", msg)

	return msg, nil
}
