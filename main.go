package main

import (
	"log"
	"net"
	"os"

	"github.com/ammardev/webserver/messages"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9000")
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
