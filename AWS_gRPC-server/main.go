package main

import (
	"aws-server/config"
	"fmt"
	"log"
	"net/http"

	"aws-server/transport"
)

func main() {
	server := transport.NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("/download", server.DownloadEndpoint)
	log.Printf("http started on addr:%v\n", config.HTTP_ADDR)
	err := http.ListenAndServe(config.HTTP_ADDR, mux)
	if err != nil {
		fmt.Println(err)
	}
}
