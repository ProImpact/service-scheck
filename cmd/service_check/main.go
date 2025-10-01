package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/ProImpact/service-check/internal/app"
	"github.com/ProImpact/service-check/internal/config"
)

var version = "v2"

var configFile = flag.String("config", "server-config.json", "Configuration file")

func main() {
	flag.Parse()
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
	slog.Info("meta", "version", version)
	application := app.NewServerApp(&cfg)
	application.Run()
}
