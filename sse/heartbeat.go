package sse

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-deeper/chunks"
	"golang.org/x/sync/errgroup"
)

func (srv *Server) BatchKeepAlive(ctx context.Context) {
	timer := time.NewTicker(20 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			go srv.batchKeepAlive(ctx)
		}
	}
}

func (srv *Server) batchKeepAlive(ctx context.Context) {
	games := srv.collectGames()
	chunkList := chunks.Split(games, 200)
	eg := new(errgroup.Group)
	for _, chunk := range chunkList {
		eg.Go(func() error {
			result, err := srv.bclient.BatchKeepAlive(ctx, chunk)
			if err != nil {
				return err
			}
			if len(result.FailedGameIDs) > 0 {
				slog.Error("batch heartbeat failed", "games_id", result.FailedGameIDs)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		slog.Error("batch heartbeat failed", "err", err)
	}
}

func (srv *Server) collectGames() (games []string) {
	srv.mux.RLock()
	defer srv.mux.RUnlock()
	for gid := range srv.games {
		games = append(games, gid)
	}
	return games
}
