package messages

import (
	"fmt"
	"net"
	"net/http"
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

func ResponseWithFile(connection net.Conn, file []byte) {
	connection.Write([]byte(SERVER_VERSION + "200 Success\r\n"))
	// MIME sniffing
	connection.Write([]byte("Content-Type: " + http.DetectContentType(file) + "\r\n"))
	connection.Write([]byte("\r\n"))
	connection.Write(file)
}
