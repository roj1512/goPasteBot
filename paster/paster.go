package paster

import (
	"bufio"
	"net"
)

func Paste(what string, to *string) error {
	conn, err := net.Dial("tcp", "ezup.dev:9999")
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(what))
	if err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	*to, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	return nil
}
