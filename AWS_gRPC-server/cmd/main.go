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

	mux := http.NewServeMux()
	mux.HandleFunc("/download", server.DownloadEndpoint)
	fmt.Printf("http started on addr:%v\n", *addr)
	err := http.ListenAndServe(*addr, mux)
	if err != nil {
		fmt.Println(err)
	}

}
