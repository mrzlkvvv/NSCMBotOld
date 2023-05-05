package bot

import "gopkg.in/telebot.v3"

func start(ctx telebot.Context) error {
    return ctx.Send("Hello, World!")
}
