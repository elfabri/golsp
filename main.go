package main

import (
	"bufio"
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
        msg := scanner.Text()
        handleMessage(logger, msg)
    }
}

func handleMessage(logger *log.Logger, msg any) {
    logger.Println(msg)
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
