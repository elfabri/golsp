package main

import (
	"bufio"
	"encoding/json"
	"golsp/lsp"
	"golsp/rpc"
	"log"
	"os"
)

func main() {
    logger := getLogger("/home/ar-fabri/lab/elfabri/golsp/log.txt")
    logger.Println("Logger started!")
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Split(rpc.Split)    // from here, scanner will wait on standar out 
                                // untill we get the message
    for scanner.Scan() {
        msg := scanner.Bytes()
        method, contents, err := rpc.DecodeMessage(msg)
        if err != nil {
            logger.Printf("Got an error: %s", err)
            continue
        }

        handleMessage(logger, method, contents)
    }
}

func handleMessage(logger *log.Logger, method string, contents []byte) {
    logger.Printf("Recived msg with method: %s", method)

    switch method {
    case "initialize":
        var request lsp.InitializeRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("Hey, we couldn't parse this: %s", err)
        }

        logger.Printf("Connected to: %s %s",
            request.Params.ClientInfo.Name,
            request.Params.ClientInfo.Version)

        // reply
        msg := lsp.NewInitializeResponse(request.ID)
        reply := rpc.EncodeMessage(msg)

        writer := os.Stdout
        writer.Write([]byte(reply))

        logger.Print("Sent the reply")

    case "textDocument/didOpen":
        var request lsp.DidOpenTextDocumentNotification
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("Hey, we couldn't parse this: %s", err)
        }

        logger.Printf("Opened: %s %s",
            request.Params.TextDocument.URI,
            request.Params.TextDocument.Text)
    }
}

// we can't print the std out, so we need to log it into a file

func getLogger(filename string) *log.Logger {
    // create a file, truncate so we get a new one every time we run this
    // and we enable write and read only, and we open it so anyone can read it with 0666
    logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
    if err != nil {
        panic("invalid filename")
    }

    return log.New(logfile, "[educationalLsp]", log.Ldate|log.Ltime|log.Lshortfile)
}
