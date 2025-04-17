package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/rpc"
	"os/signal"
	"syscall"

	"github.com/bketelsen/inventory"
	"github.com/bketelsen/inventory/routes"
	"github.com/bketelsen/inventory/service"
	"github.com/bketelsen/inventory/storage"
	"github.com/bketelsen/toolbox/cobra"
	"github.com/coreos/go-systemd/daemon"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

// NewServerCommand creates a new server command
func NewServerCommand(config *viper.Viper) *cobra.Command {
	// Define our command
	serverCmd := &cobra.Command{
		Use:           "server",
		SilenceErrors: true,
		Aliases:       []string{"serve"},
		Short:         "starts the RPC and HTTP servers",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Create a new memory storage
			memStorage := storage.NewMemoryStorage()

			// Create a new inventory server with the storage
			server := service.NewInventoryServer(memStorage)
			err := rpc.Register(server) // Register the Inventory service
			if err != nil {
				log.Printf("Error registering RPC server: %v", err)
				return err
			}

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

			return run(cmd.Context(), cmd.Logger, strconv.Itoa(config.GetInt("server.http-port")), memStorage)
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

func run(ctx context.Context, logger *slog.Logger, httpPort string, storage inventory.Storage) error {
	sctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg, ctx := errgroup.WithContext(sctx)

	eg.Go(func() error {
		router := chi.NewMux()

		router.Use(
			middleware.Recoverer,
		)

		router.Handle("/static/*", http.StripPrefix("/static/", inventory.Static(logger)))

		if err := routes.SetupRoutes(ctx, router, storage); err != nil {
			return fmt.Errorf("error setting up routes: %w", err)
		}

		srv := &http.Server{
			Addr:    "0.0.0.0:" + httpPort,
			Handler: router,
		}

		go func() {
			<-ctx.Done()
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Fatalf("error during shutdown: %v", err)
			}
		}()
		logger.Info("HTTP server started", "port", httpPort)
		return srv.ListenAndServe()
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
