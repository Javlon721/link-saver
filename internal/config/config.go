package config

type Config struct {
	// Bot
	TELEGRAM_TOKEN string
	BOT_NAME       string

	// Database
	POSTGRES_DB       string
	POSTGRES_PASSWORD string
	POSTGRES_USER     string
	POSTGRES_HOST     string
	POSTGRES_PORT     uint16
}
