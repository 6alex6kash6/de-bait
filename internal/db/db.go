package db

import (
	"context"
	"log/slog"

	"github.com/de-bait/ent"
	"go.uber.org/fx"
)

func NewEntClient(lc fx.Lifecycle) *ent.Client {
	client, err := ent.Open("postgres", "postgresql://postgres:3222820schism@db.omhvdsjuyetldestuyqx.supabase.co:5432/postgres")
	if err != nil {
		slog.Error("failed opening connection to postgres: %v", err)
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Close()
		},
	})
	return client
}
