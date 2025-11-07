package main

import (
	"github.com/sung2708/shorten_url/internal/config"
	"github.com/sung2708/shorten_url/internal/database"
	"github.com/sung2708/shorten_url/internal/router"
)

func main() {
	cfg := config.NewConfigFromEnv()

	db := database.InitPostgres(cfg.PostgresDSN)
	//rdb := database.InitRedis(cfg.RedisHost)

	r := router.Setup(cfg, db)
	err := r.Run(":" + cfg.PORT)
	if err != nil {
		return
	}
}
