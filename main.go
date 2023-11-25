package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/ammardev/webserver/messages"
)

func main() {
	host := flag.String("host", "localhost", "The host to serve HTTP on")
	port := flag.Int("port", 9000, "The port we serve HTTP on")
	debug := flag.Bool("debug", false, "Show debug logging")

	flag.Parse()

	var loggingLevel slog.Level

	if *debug {
		loggingLevel = slog.LevelDebug
	} else {
		loggingLevel = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggingLevel,
	})))

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		log.Fatalln(err)
	}

	slog.Info(
		"Start serving HTTP",
		"host", *host,
		"port", *port,
	)

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
	slog.Debug("Handling new connection", "addr", connection.RemoteAddr())
	defer func() {
		slog.Debug("Closing connection", "addr", connection.RemoteAddr())
		connection.Close()
	}()

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
