package rpc_test

import (
	"golsp/rpc"
	"testing"
)

type EncodingExample struct {
    Testing bool
}

func TestEncode(t *testing.T) {
    // Content-Length: <NUMBER>\r\n\r\n{\"testing\":true}
    // {"testing:true"} -> 16 characters -> 16 bytes
    expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
    actual := rpc.EncodeMessage(EncodingExample{Testing:true})

    if expected != actual {
        t.Fatalf("Expected: %s, Actual: %s", expected, actual)
    }
}

func TestDecode(t *testing.T) {
    incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
    method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
    contentLength := len(content)
    if err != nil {
        t.Fatal(err)
    }

    if contentLength != 15 {
        t.Fatalf("Expected: 15, got: %d", contentLength)
    }

    if method != "hi" {
        t.Fatalf("Expected: 'hi', got: %s", method)
    }
}
