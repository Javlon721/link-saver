package main

import (
	"github.com/Javlon721/link-saver/internal/config"
	linksaver "github.com/Javlon721/link-saver/internal/linkSaver"
	"github.com/joho/godotenv"
)

func main() {
	envs, err := godotenv.Read(".env")

	if err != nil {
		panic(err)
	}

	config := &config.Config{
		TELEGRAM_TOKEN: envs["TELEGRAM_TOKEN"],
		BOT_NAME:       envs["BOT_NAME"],
	}

	app, err := linksaver.New(config)

	if err != nil {
		panic(err)
	}

	app.Start()
}
