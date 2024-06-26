package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
    content, err := json.Marshal(msg) // encoding the msg

    // if error, stop 
    if err != nil {
        panic(err)
    }

    return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
} 

type BaseMessage struct {
    Method string `json:"method"`
}

func DecodeMessage(msg []byte) (string, []byte, error) {
    header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})

    if !found {
        return "", nil, errors.New("Did not find separator")
    }

    // Content-Length: <number>
    // we get the length in bytes, after 'Content-Length: '
    contentLengthBytes := header[len("Content-Length: "):]
    contentLength, err := strconv.Atoi(string(contentLengthBytes))

    if err != nil {
        return "", nil, err
    }

    var baseMessage BaseMessage
    if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
        return "", nil, err
    }

    return baseMessage.Method, content[:contentLength], err
}

// type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
// we are not usign atEOF (at End Of File)

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
    header, content, found := bytes.Cut(data, []byte{'\r', '\n', '\r', '\n'})

    if !found {
        return 0, nil, nil
    }

    // Content-Length: <number>
    // we get the length in bytes, after 'Content-Length: '
    contentLengthBytes := header[len("Content-Length: "):]
    contentLength, err := strconv.Atoi(string(contentLengthBytes))

    // invalid contentLength
    if err != nil {
        return 0, nil, err
    }

    // we can keep reading
    if len(content) < contentLength {
        return 0, nil, nil
    }

    // \r\n\r\n -> 4 bytes
    totalLength := len(header) + 4 + contentLength

    return totalLength, data[:totalLength], nil
}

