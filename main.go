package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ammardev/webserver/messages"
)

func main() {
	host := flag.String("host", "localhost", "The host to serve HTTP on")
	port := flag.Int("port", 9000, "The port we serve HTTP on")

	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
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
		messages.ResponseWithError(connection, err)
		return
	}

	file, err := getFileContentsByURI(request.URI)
	if err != nil {
		messages.ResponseWithError(connection, err)
		return
	}

	messages.ResponseWithFile(connection, file)
}

func getFileContentsByURI(uri string) ([]byte, messages.HttpError) {
	fileName := "index.html"

	if uri != "/" {
		fileName = uri
	}

	file, err := os.ReadFile("./public/" + fileName)
	if err != nil {
		return nil, messages.NotFoundErr{}
	}

	return []byte(file), nil
}
