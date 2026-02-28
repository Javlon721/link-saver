package linksaver

import (
	"fmt"
	"time"

	"github.com/Javlon721/link-saver/internal/config"
	tele "gopkg.in/telebot.v4"
)

type Registable interface {
	RegisterHandlers(*tele.Bot)
}

type App struct {
	Bot    *tele.Bot
	Config *config.Config
}

func New(config *config.Config) (*App, error) {
	app := &App{Config: config}

	pref := tele.Settings{
		Token:  config.TELEGRAM_TOKEN,
		Poller: &tele.LongPoller{Timeout: 15 * time.Second},
	}

	bot, err := tele.NewBot(pref)

	if err != nil {
		return nil, err
	}

	app.Bot = bot

	return app, nil
}

func (app App) Start() {

	fmt.Printf("%s started to work\n", app.Config.BOT_NAME)

	app.Bot.Start()
}

func (app App) RegisterHandler(h Registable) {
	h.RegisterHandlers(app.Bot)
}
