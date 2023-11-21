package main

import (
	"log"
	"net"
	"os"

	"github.com/ammardev/webserver/messages"
)

func main() {
    listener, err := net.Listen("tcp", "0.0.0.0:9000");
    if err != nil {
        log.Fatalln(err)
    }

    defer listener.Close()

    for {
        connection, err := listener.Accept()
        if err != nil {
            log.Fatalln(err)
        }

        go handleConnection(connection)
    }
}


func handleConnection(connection net.Conn) {
    defer connection.Close()

    request, err := messages.NewRequest(connection)
    if err != nil {
        connection.Write([]byte(err.Error() + "\r\n"))
        return
    }

    fileName := "index.html"

    if request.URI != "/" {
        fileName = request.URI
    }

    file, err := os.ReadFile("./public/" + fileName)
    if err != nil {
        connection.Write([]byte("HTTP/1.0 404 File Not Found\r\n"))
        return
    }

    connection.Write([]byte("HTTP/1.0 200 Success\r\n\r\n"))
    connection.Write(file)
}
