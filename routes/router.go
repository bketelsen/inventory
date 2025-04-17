// Package routes defines the routes for the web server
package routes

import (
	"context"

	"github.com/bketelsen/inventory"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(_ context.Context, router chi.Router, storage inventory.Storage) (err error) {
	setupIndexRoute(router, storage)

	return nil
}
