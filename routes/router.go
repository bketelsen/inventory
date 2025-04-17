package routes

import (
	"context"
	"errors"
	"fmt"

	"github.com/bketelsen/inventory"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(ctx context.Context, router chi.Router, storage inventory.Storage) (err error) {
	if err := errors.Join(
		setupIndexRoute(router, storage),
	); err != nil {
		return fmt.Errorf("error setting up routes: %w", err)
	}

	return nil
}
