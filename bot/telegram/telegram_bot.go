package telegram

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/indiependente/motionbot/bot"
	"github.com/pkg/errors"
)

type UserInfo struct {
	chatID int64
}
type Bot struct {
	bot          *tgbotapi.BotAPI
	allowedUsers map[string]UserInfo
	updates      <-chan tgbotapi.Update
}

func NewBot(tok string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(tok)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create bot")
	}
	return &Bot{bot, make(map[string]UserInfo), nil}, nil
}

func (b *Bot) Setup(bc bot.BotConfig) error {
	for _, usrid := range bc.AllowedUsers {
		uid := strings.Split(usrid, ":")
		cid, err := strconv.Atoi(uid[1])
		if err != nil {
			return errors.Wrapf(err, "Could not convert user %s chat id %s", uid[0], uid[1])
		}
		b.allowedUsers[uid[0]] = UserInfo{chatID: int64(cid)}
	}
	b.bot.Debug = true
	return nil
}
func (b *Bot) Start() error {
	var err error
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	b.updates, err = b.bot.GetUpdatesChan(u)
	if err != nil {
		return errors.Wrap(err, "Could not create updates channel")
	}
	// start an updates goroutine handler here
	go b.updateHandler(b.updates)
	return nil
}

func (b *Bot) Send(m bot.Message) error {
	for usr, usrinfo := range b.allowedUsers {
		switch m.Format {
		case bot.TEXT:
			msg := tgbotapi.NewMessage(usrinfo.chatID, string(m.Data))
			_, err := b.bot.Send(msg)
			if err != nil {
				return errors.Wrapf(err, "Could not send text message to %s", usr)
			}
		case bot.IMAGE:
			msg := tgbotapi.NewPhotoUpload(usrinfo.chatID, string(m.Data))
			msg.Caption = string(m.Data)
			_, err := b.bot.Send(msg)
			if err != nil {
				return errors.Wrapf(err, "Could not send photo message to %s", usr)
			}
		case bot.VIDEO:
			msg := tgbotapi.NewVideoNoteUpload(usrinfo.chatID, 240, string(m.Data))
			msg.Duration = 10
			_, err := b.bot.Send(msg)
			if err != nil {
				return errors.Wrapf(err, "Could not send video note message to %s", usr)
			}
		}

	}

	return nil
}

func (b *Bot) updateHandler(updates <-chan tgbotapi.Update) {
	for u := range updates {
		if u.Message == nil { // ignore any non-Message Updates
			continue
		}
		switch u.Message.Text {
		case "/subscribe":
			if _, ok := b.allowedUsers[u.Message.From.UserName]; ok {
				b.allowedUsers[u.Message.From.UserName] = UserInfo{chatID: u.Message.Chat.ID}
			}
		}
		log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)
	}
}
