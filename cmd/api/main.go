package main

import (
	"github.com/sung2708/shorten_url/internal/config"
	"github.com/sung2708/shorten_url/internal/database"
	"github.com/sung2708/shorten_url/internal/router"
)

func main() {
	cfg := config.NewConfigFromEnv()

	db := database.InitPostgres(cfg.PostgresDSN)
	rdb := database.InitRedis(cfg.RedisHost, cfg.RedisUser, cfg.RedisPass)

	r := router.Setup(cfg, db, rdb)
	err := r.Run(":" + cfg.PORT)
	if err != nil {
		return
	}
}
