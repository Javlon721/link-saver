package main

import (
	"context"
	"strconv"

	"github.com/Javlon721/link-saver/internal/config"
	"github.com/Javlon721/link-saver/internal/db"
	"github.com/Javlon721/link-saver/internal/handlers"
	linksaver "github.com/Javlon721/link-saver/internal/linkSaver"
	"github.com/Javlon721/link-saver/internal/middleware"
	"github.com/Javlon721/link-saver/internal/services"
	"github.com/joho/godotenv"
)

var (
	userTable = "users"
	linkTable = "links"
)

func main() {
	envs, err := godotenv.Read(".env")

	if err != nil {
		panic(err)
	}

	postgresPort, err := strconv.ParseUint(envs["POSTGRES_PORT"], 10, 16)

	config := &config.Config{
		TELEGRAM_TOKEN: envs["TELEGRAM_TOKEN"],
		BOT_NAME:       envs["BOT_NAME"],

		POSTGRES_DB:       envs["POSTGRES_DB"],
		POSTGRES_PASSWORD: envs["POSTGRES_PASSWORD"],
		POSTGRES_USER:     envs["POSTGRES_USER"],
		POSTGRES_HOST:     envs["POSTGRES_HOST"],
		POSTGRES_PORT:     uint16(postgresPort),
	}

	postgreConn, err := db.NewPostgreConn(config)

	if err != nil {
		panic(err)
	}

	defer postgreConn.Close(context.Background())

	userStore := db.NewPostgresUserStore(postgreConn, userTable)
	linkStore := db.NewPostgreLinkStore(postgreConn, linkTable)

	linkService := services.NewLinkService(linkStore)
	userService := services.NewUserService(userStore)

	app, err := linksaver.New(config)

	if err != nil {
		panic(err)
	}

	mainHandler := handlers.NewMainHandler()
	userHandler := handlers.NewUserHandler(userService, linkService, postgreConn)
	linkHandler := handlers.NewLinkHandler(linkService)

	authMiddleware := middleware.AuthorizeUser(userService)

	// global use endpoints
	app.Bot.Handle("/start", mainHandler.HelpDeskHandler)
	app.Bot.Handle("/help", mainHandler.HelpDeskHandler)

	// register new user
	app.Bot.Handle("/register", userHandler.RegisterUser)

	// user crud
	app.Bot.Handle("/me", userHandler.GetUser, authMiddleware)
	app.Bot.Handle("/stop", userHandler.DeleteUser, authMiddleware)

	// link crud
	app.Bot.Handle("/link", linkHandler.RegisterLink, authMiddleware)
	app.Bot.Handle("/links", linkHandler.GetAll, authMiddleware)
	app.Bot.Handle("/linksBtns", linkHandler.GetAllWithBtns, authMiddleware)

	//link crud callbacks
	app.RegisterCallbacks(linkHandler.GetCallbacks())
	app.ListenCallbacks(authMiddleware)

	app.Start()
}
