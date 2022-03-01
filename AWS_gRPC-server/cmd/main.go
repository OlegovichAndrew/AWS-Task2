package main

import (
	"flag"
	"fmt"
	"net/http"

	"aws-server/transport"
)

var (
	addr = flag.String("addr", "localhost:4445", "The server address")
)

func main() {
	flag.Parse()
	server := transport.NewServer()

	http.HandleFunc("/download", server.DownloadEndpoint)
	fmt.Println("http started")
	http.ListenAndServe(*addr, nil)

}
