/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	netrpc "net/rpc"
	"os"

	"github.com/ProImpact/service-check/internal/rpc"
	"github.com/spf13/cobra"
)

var serviceName string

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "See the logs of a running process",
	Run: func(cmd *cobra.Command, args []string) {
		if serviceName == "" {
			slog.Error("Service name arg is unset")
			cmd.Help()
			os.Exit(1)
		}
		client, err := netrpc.DialHTTP("tcp", serverPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer client.Close()
		rpcArgs := &rpc.ServiceNameArgs{
			ServiceName: serviceName,
		}
		reply := &rpc.ServiceLogsReply{}
		err = client.Call("Server.Logs", rpcArgs, reply)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		fmt.Println(reply.Data)
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().StringVar(
		&serviceName,
		"name",
		"",
		"Service name to kill",
	)
}
