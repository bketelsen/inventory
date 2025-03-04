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
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/bketelsen/inventory/service"
	"github.com/bketelsen/inventory/storage"
	"github.com/bketelsen/inventory/types"
	"github.com/bketelsen/inventory/web"
	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve starts the RPC and HTTP servers",

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := types.ReadConfig()
		if err != nil {
			log.Println("Error reading config:", err)
			return err
		}
		log.Println("Starting inventory server")
		log.Printf("HTTP Port: %d", cfg.HTTPPort)
		log.Printf("RPC Port: %d", cfg.RPCPort)
		// Create a new memory storage
		memStorage := storage.NewMemoryStorage(cfg)

		// Create a new inventory server with the storage
		server := service.NewInventoryServer(cfg, memStorage)
		rpc.Register(server) // Register the Inventory service

		sent, e := daemon.SdNotify(false, "READY=1")
		if e != nil {
			log.Fatal(e)
		}
		if !sent {
			log.Printf("SystemD notify NOT sent\n")
		}

		listener, err := net.Listen("tcp", ":9999")
		if err != nil {
			log.Printf("Error listening: %v", err)
			return err
		}
		defer listener.Close()

		log.Println("RPC Server is listening on port 9999...")
		go func() {
			rpc.Accept(listener) // Accept incoming RPC connections
		}()
		// Use a template that doesn't take parameters.
		//	http.Handle("/", templ.Handler(home()))

		http.Handle("/static/", http.FileServer(http.FS(web.Static)))

		http.Handle("/", web.NewInventoryHandler(memStorage))
		log.Println("http server listening on :8000")
		if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
			log.Printf("error listening: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")
	serveCmd.PersistentFlags().IntP("http-port", "w", 8000, "HTTP Port")
	viper.BindPFlag("http_port", serveCmd.PersistentFlags().Lookup("http-port"))
	serveCmd.PersistentFlags().IntP("rpc-port", "r", 9999, "RPC Port")
	viper.BindPFlag("rpc_port", serveCmd.PersistentFlags().Lookup("rpc-port"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
