package telegram

import (
	"database/sql"
	"doMassageBot/internal/config"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	uId  int
	mId  int
	conf *config.Configuration
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		bot: bot,
	}
}

func (b *Bot) Start(db *sql.DB, conf *config.Configuration) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	b.conf = conf

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message != nil {
			var err error
			if update.Message.From.ID == conf.AdminID {
				err = b.adminHandler(update.Message, db)
			} else {
				if update.Message.IsCommand() {
					err = b.handleCommand(update.Message, db)
				} else {
					err = b.handleText(update.Message, db)
				}
			}

			if err != nil {
				return err
			}

			continue
		}

		if update.CallbackQuery != nil {
			if err := b.handleCallbackQuery(update.CallbackQuery, db); err != nil {
				return err
			}
		}
	}
	return nil
}
