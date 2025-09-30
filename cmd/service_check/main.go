package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/ProImpact/passboult/internal/app"
	"github.com/ProImpact/passboult/internal/config"
)

var configFile = flag.String("config", "server-config.json", "Configuration file")

func main() {
	f, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fileData, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var cfg config.ServerConfiguration
	err = json.Unmarshal(fileData, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	application := app.NewServerApp(&cfg)
	application.Run()
}
