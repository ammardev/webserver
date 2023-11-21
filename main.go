package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"os"
	"strings"
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

    connectionReader := bufio.NewReader(connection)

    requestLine, _, err := connectionReader.ReadLine()
    if err != nil {
        log.Fatalln(err)
    }


    if !bytes.Equal(requestLine[:3], []byte("GET")) {
        connection.Write([]byte("HTTP/1.0 501 Only GET is supported\r\n"))
        return
    }

    var pathBuilder strings.Builder

    for _, byte := range requestLine[4:] {
        if byte == ' ' {
            break
        }
        pathBuilder.WriteByte(byte)
    }

    path := pathBuilder.String()

    if path == "/" {
        path = "index.html"
    }

    file, err := os.ReadFile("./public/" + path)
    if err != nil {
        connection.Write([]byte("HTTP/1.0 404 File Not Found\r\n"))
        return
    }

    connection.Write([]byte("HTTP/1.0 200 Success\r\n\r\n"))
    connection.Write(file)
}
