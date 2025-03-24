package common

import (
	"bufio"
	"net"
	"io"
	"encoding/binary"
)



func send_message(conn net.Conn, id int, bets []ClientData) (string, error) {
	var messages []string
	for _, bet := range bets {
		message := fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s", id,bet.Nombre, bet.Apellido, bet.Documento, bet.Nacimiento, bet.Numero)
		messages = append(messages, message)
	}
	
	batchMessage := strings.Join(messages, ">")

	
	len := len(batchMessage)
	log.Infof("action: len_message | result: success | len: %v", len)


	binary.Write(conn, binary.BigEndian, uint16(len))
	io.WriteString(conn, batchMessage)

	// Leemos la respuesta del servidor
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Errorf("action: receive_message | result: fail | error: %v", err)
		return "", err
	}

	log.Infof("action: receive_message | result: success | server_response: %v", msg)
	return msg, err
}


func send_message(conn net.Conn, id int, bets []ClientData, maxBatchSize int) (string, error) {
	var messages []string
	var batchMessage string

	for _, bet := range bets {
		message := fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s", id,bet.Nombre, bet.Apellido, bet.Documento, bet.Nacimiento, bet.Numero)
		messages = append(messages, message)
	}

	// Ahora, enviamos el batch en partes si el tamaño total excede el límite
	for len(messages) > 0 {
		batchMessage = ""
		for len(messages) > 0 {
			// Concatenamos la apuesta en el batch, usando '>' como delimitador entre mensajes
			batchMessage += messages[0] + ">"

			if len(batchMessage) > maxBatchSize {
				// Si el tamaño supera el limite, eliminamos el último '>' y enviamos el batch
				batchMessage = batchMessage[:len(batchMessage)-1]
				break
			}
			messages = messages[1:]
		}

		len := len(batchMessage)
		log.Infof("action: len_message | result: success | len: %v", len)
		binary.Write(conn, binary.BigEndian, uint16(len))
		io.WriteString(conn, batchMessage)

		// Leemos la respuesta del servidor
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Errorf("action: receive_message | result: fail | error: %v", err)
			return "", err
		}

		log.Infof("action: receive_message | result: success | server_response: %v", msg)

	}

	return msg, err
}
