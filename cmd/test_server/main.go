package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

var port = flag.Int("port", 3000, "Test server port")

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Ok"))
	})
	addrs := fmt.Sprintf(":%d", *port)
	slog.Info("server started", "port", *port)
	slog.Info("new hello from web server")
	log.Fatal(http.ListenAndServe(addrs, router))
}
