/*
Copyright © 2025 Brian Ketelsen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/rpc"

	"github.com/bketelsen/inventory/service"
	"github.com/bketelsen/inventory/storage"
	"github.com/bketelsen/inventory/types"
	"github.com/bketelsen/inventory/web"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/coreos/go-systemd/v22/daemon"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve starts the RPC and HTTP servers",

	RunE: func(cmd *cobra.Command, args []string) error {
		slog.SetDefault(cmd.Logger)

		cfg, err := types.ViperToStruct(cmd.GlobalConfig())
		if err != nil {
			log.Println("Error reading config:", err)
			return err
		}
		cmd.Logger.Info("Starting inventory server", "httpport", cfg.HTTPPort, "rpcport", cfg.RPCPort)

		// set up logging since this is a long running process

		// Create a new memory storage
		memStorage := storage.NewMemoryStorage()

		// Create a new inventory server with the storage
		server := service.NewInventoryServer(memStorage)
		rpc.Register(server) // Register the Inventory service

		sent, e := daemon.SdNotify(false, "READY=1")
		if e != nil {
			log.Fatal(e)
		}
		if !sent {
			cmd.Logger.Warn("SystemD notify NOT sent")
		}

		rpcPort := cmd.Config().GetInt("rpc-port")
		rpcStr := fmt.Sprintf(":%d", rpcPort)
		listener, err := net.Listen("tcp", rpcStr)
		if err != nil {
			log.Printf("Error listening: %v", err)
			return err
		}
		defer listener.Close()

		cmd.Logger.Info("RPC Server is listening", "port", rpcPort)
		go func() {
			rpc.Accept(listener) // Accept incoming RPC connections
		}()

		http.Handle("/static/", http.FileServer(http.FS(web.Static)))

		http.Handle("/", web.NewInventoryHandler(memStorage))

		httpPort := cmd.Config().GetInt("http-port")
		httpStr := fmt.Sprintf("0.0.0.0:%d", httpPort)
		cmd.Logger.Info("http server listening", "port", httpPort)
		if err := http.ListenAndServe(httpStr, nil); err != nil {
			cmd.Logger.Error(fmt.Sprintf("error listening: %v", err))
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().IntP("http-port", "w", 8000, "HTTP port")
	serveCmd.PersistentFlags().IntP("rpc-port", "r", 9999, "RPC port")

}
