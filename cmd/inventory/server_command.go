package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/bketelsen/inventory/service"
	"github.com/bketelsen/inventory/storage"
	"github.com/bketelsen/inventory/web"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/coreos/go-systemd/daemon"
	"github.com/spf13/viper"
)

// Build the cobra command that handles our command line tool.
func NewServerCommand(config *viper.Viper) *cobra.Command {

	// Define our command
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "starts the RPC and HTTP servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a new memory storage
			memStorage := storage.NewMemoryStorage()

			// Create a new inventory server with the storage
			server := service.NewInventoryServer(memStorage)
			rpc.Register(server) // Register the Inventory service

			sent, _ := daemon.SdNotify(false, "READY=1")

			if !sent {
				cmd.Logger.Warn("Not running under systemd, not sending notification")
			}
			listen := config.GetString("server.listen")

			rpcPort := config.GetInt("server.rpc-port")
			rpcStr := fmt.Sprintf("%s:%d", listen, rpcPort)
			listener, err := net.Listen("tcp", rpcStr)
			if err != nil {
				log.Printf("Error listening: %v", err)
				return err
			}
			defer listener.Close()

			cmd.Logger.Info("RPC server started", "address", listen, "port", rpcPort)
			go func() {
				rpc.Accept(listener) // Accept incoming RPC connections
			}()

			http.Handle("/static/", http.FileServer(http.FS(web.Static)))

			http.Handle("/", web.NewInventoryHandler(memStorage))

			httpPort := config.GetInt("server.http-port")
			httpStr := fmt.Sprintf("%s:%d", listen, httpPort)
			cmd.Logger.Info("HTTP server started", "address", listen, "port", httpPort)
			if err := http.ListenAndServe(httpStr, nil); err != nil {
				cmd.Logger.Error(fmt.Sprintf("error listening: %v", err))
				return err
			}
			return nil
		},
	}

	// define flags.
	serverCmd.Flags().StringP("listen", "l", "0.0.0.0", "Address to listen on")
	serverCmd.Flags().IntP("http-port", "w", 8000, "HTTP Port to listen on")
	serverCmd.Flags().IntP("rpc-port", "r", 9999, "RPC Port to listen on")
	serverCmd.PreRunE = func(cmd *cobra.Command, _ []string) error {
		// bind the flags to viper at "parent.child" for nested keys
		// which maps to the config file
		// server:
		//   listen: 0.0.0.0
		// etc.
		_ = config.BindPFlag("server.listen", cmd.Flags().Lookup("listen"))
		_ = config.BindPFlag("server.http-port", cmd.Flags().Lookup("http-port"))
		_ = config.BindPFlag("server.rpc-port", cmd.Flags().Lookup("rpc-port"))
		return nil
	}
	return serverCmd
}
