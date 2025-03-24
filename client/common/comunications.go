package common

import (
	"bufio"
	"net"
	"io"
	"encoding/binary"
)




func send_message(conn net.Conn, message string) (string, error) {


	len := len(message)
	log.Infof("action: len_message | result: success | len: %v", len)
	binary.Write(conn, binary.BigEndian, uint16(len))
	io.WriteString(conn, message)
	msg, err := bufio.NewReader(conn).ReadString('\n')
	return msg, err
}
