package linksaver

import (
	"fmt"
	"log/slog"
	"maps"
	"strings"
	"time"

	"github.com/Javlon721/link-saver/internal/config"
	tele "gopkg.in/telebot.v4"
)

type App struct {
	Bot       *tele.Bot
	Config    *config.Config
	Callbacks map[string]tele.HandlerFunc
}

func New(config *config.Config) (*App, error) {
	app := &App{Config: config, Callbacks: map[string]tele.HandlerFunc{}}

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

func (app App) RegisterCallbacks(callbacks map[string]tele.HandlerFunc) {
	maps.Copy(app.Callbacks, callbacks)
}

func (app App) ListenCallbacks(mids ...tele.MiddlewareFunc) {
	const startWith = "\f"

	app.Bot.Handle(tele.OnCallback, func(c tele.Context) error {
		payload := strings.TrimLeft(c.Callback().Data, startWith)

		source := strings.SplitN(payload, "|", 2)

		cb, ok := app.Callbacks[source[0]]

		if !ok {
			return nil
		}

		err := cb(c)

		if err != nil {
			slog.Error("tele.OnCallback", "err", err)
		}

		return nil
	}, mids...)
}
