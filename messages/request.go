package messages

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"
)

type HttpHeaders map[string]string

type Request struct {
	Method  string
	URI     string
	Version string
	Headers HttpHeaders
	Body    []byte
}

func NewRequest(connection net.Conn) (*Request, error) {
	reader := newRequestReader(connection)

	request := Request{}

	err := request.setMethod(reader.ReadToSP())
	if err != nil {
		return nil, err
	}

	request.URI = reader.ReadToSP()

	err = request.setVersion(reader.ReadToCRLF())
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (req *Request) setMethod(method string) error {
	switch method {
	case "GET":
		req.Method = method
	default:
		return errors.New("HTTP/1.0 501 Only GET is supported")
	}

	return nil
}

func (req *Request) setVersion(version string) error {
	switch version {
	case "HTTP/1.0", "HTTP/1.1":
		req.Version = version
	default:
		return errors.New("HTTP/1.0 505 HTTP Version Not Supported")
	}

	return nil
}

// TODO: Move connection related types to another package.
type RequestReader struct {
	bufioReader *bufio.Reader
}

func newRequestReader(connection net.Conn) RequestReader {
	return RequestReader{
		bufioReader: bufio.NewReader(connection),
	}
}

func (reader *RequestReader) ReadToSP() string {
	token, err := reader.bufioReader.ReadString(' ')
	if err != nil {
		log.Println(err)
		return ""
	}

	return strings.TrimSpace(token)
}

func (reader *RequestReader) ReadToCRLF() string {
	token, err := reader.bufioReader.ReadString('\r')
	if err != nil {
		log.Println(err)
		return ""
	}

	nextByte, err := reader.bufioReader.ReadByte()
	if err != nil {
		log.Println(err)
		return ""
	}

	if nextByte != '\n' {
		log.Println("Parsing Error. TODO: Add proper parsing")
		return ""
	}

	return strings.TrimSpace(token)
}
