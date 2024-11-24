package server

import (
	"context"

	"github.com/ScooballyD/gsource-lib/scrapers"
)

// work in prog
func (cfg *apiConfig) CollectGames() ([]scrapers.Game, error) {
	_, err := cfg.db.GetGames(context.Background())
	return nil, err

}
