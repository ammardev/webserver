package messages

import (
	"fmt"
	"net"
)

const SERVER_VERSION = "HTTP/1.1"

func ResponseWithError(connection net.Conn, err HttpError) {
	connection.Write([]byte(fmt.Sprintf(
		"%s %d %s\r\n",
		SERVER_VERSION,
		err.Status(),
		err.Error(),
	)))
}

func ResponseWithHtmlFile(connection net.Conn, file []byte) {
	connection.Write([]byte(SERVER_VERSION + "200 Success\r\n"))
	connection.Write([]byte("Content-Type: text/html\r\n"))
	connection.Write([]byte("\r\n"))
	connection.Write(file)
}
