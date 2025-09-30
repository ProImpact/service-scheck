/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	netrpc "net/rpc"
	"os"

	"github.com/ProImpact/passboult/internal/rpc"
	"github.com/ProImpact/passboult/pkg"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all the services managed by the service manager",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := netrpc.DialHTTP("tcp", serverPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		defer client.Close()
		rpcArgs := struct{}{}
		var reply rpc.ServiceListReply
		err = client.Call("Server.GetAllServices", rpcArgs, &reply)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		if len(reply.Services) == 0 {
			slog.Info("No service running")
			return
		}
		pkg.MustPrint(pkg.Event{
			Data: reply.Services,
		})
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
