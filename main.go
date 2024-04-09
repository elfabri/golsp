package main

import (
	"bufio"
	"encoding/json"
	"golsp/analysis"
	"golsp/lsp"
	"golsp/rpc"
	"io"
	"log"
	"os"
)

func main() {
    logger := getLogger("/home/ar-fabri/lab/elfabri/golsp/log.txt")
    logger.Println("Logger started!")
    scanner := bufio.NewScanner(os.Stdin)

    scanner.Split(rpc.Split)    // from here, scanner will wait on standar out 
                                // untill we get the message

    state := analysis.NewState()
    writer := os.Stdout

    for scanner.Scan() {
        msg := scanner.Bytes()
        method, contents, err := rpc.DecodeMessage(msg)
        if err != nil {
            logger.Printf("Got an error: %s", err)
            continue
        }

        handleMessage(logger, writer, state, method, contents)
    }
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, contents []byte) {
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
        writeResponse(writer, msg)


        logger.Print("Sent the reply")

    case "textDocument/didOpen":
        var request lsp.DidOpenTextDocumentNotification
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/didOpen: %s", err)
            return
        }

        logger.Printf("Opened: %s", request.Params.TextDocument.URI,)

        // sync the state of the document (logger(server)) with the state of
        // the file (editor(client)) at the opening of that file
        state.OpenDocument(
            request.Params.TextDocument.URI,
            request.Params.TextDocument.Text)

    case "textDocument/didChange":
        var request lsp.TextDocumentDidChangeNotification
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/didChange %s", err)
            return
        }

        logger.Printf("Changed: %s", request.Params.TextDocument.URI,)

        for _, change := range request.Params.ContentChanges {
            state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
        }

    case "textDocument/hover":
        var request lsp.HoverRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/hover %s", err)
            return
        }

        // Create Response and write it back
        response := state.Hover(
            request.ID,
            request.Params.TextDocument.URI,
            request.Params.Position)

        writeResponse(writer, response)

    case "textDocument/definition":
        var request lsp.DefinitionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/definition %s", err)
            return
        }

        // Create Response and write it back
        response := state.Definition(
            request.ID,
            request.Params.TextDocument.URI,
            request.Params.Position)

        writeResponse(writer, response)

    case "textDocument/codeAction":
        var request lsp.CodeActionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/codeAction: %s", err)
            return
        }

        // create response
        response := state.TextDocumentCodeAction(
            request.ID,
            request.Params.TextDocument.URI)

        // write it back
        writeResponse(writer, response)

    case "textDocument/completion":
        var request lsp.CompletionRequest
        if err := json.Unmarshal(contents, &request); err != nil {
            logger.Printf("textDocument/completion: %s", err)
            return
        }

        // create response
        response := state.TextDocumentCompletion(
            request.ID,
            request.Params.TextDocument.URI)

        // write it back
        writeResponse(writer, response)
    }

}

func writeResponse(writer io.Writer, msg any) {
    reply := rpc.EncodeMessage(msg)
    writer.Write([]byte(reply))
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
