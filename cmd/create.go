/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"io"
	"log/slog"
	netrpc "net/rpc"
	"os"
	"strings"

	"github.com/ProImpact/service-check/internal/rpc"
	"github.com/ProImpact/service-check/pkg/model"
	"github.com/spf13/cobra"
)

var serviceDefinitionPath string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Registrer a new service process",
	Long:  `Regitrer a new service by the process manager`,
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Open(serviceDefinitionPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer f.Close()
		fileData, err := io.ReadAll(f)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		var m model.ServiceCreate
		err = json.Unmarshal(fileData, &m)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		client, err := netrpc.DialHTTP("tcp", serverPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer client.Close()
		m.Check.Type = strings.ToUpper(m.Check.Type)
		rpcArgs := &rpc.ServiceArgs{
			Service: model.ServiceCreate{
				ServiceName: m.ServiceName,
				Check:       m.Check,
				PingTime:    m.PingTime,
				Command:     m.Command,
			},
		}
		reply := false
		err = client.Call("Server.CreateService", rpcArgs, &reply)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		if reply {
			slog.Info("Service created", "name", rpcArgs.Service.ServiceName)
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&serviceDefinitionPath, "path", "service.json", "Service definition path")
}
