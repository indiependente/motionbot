package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/indiependente/motionbot/sensor/camera"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/indiependente/motionbot/bot"
	"github.com/pkg/errors"
)

const (
	helloMessageTemplate = "Welcome %s! 😃"
)

type UserInfo struct {
	username  string
	firstName string
	lastName  string
	chatID    int64
}
type Bot struct {
	bot          *tgbotapi.BotAPI
	allowedUsers map[int]UserInfo
	updates      <-chan tgbotapi.Update
	camera       camera.Camera
}

func NewBot(tok string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(tok)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create bot")
	}
	return &Bot{bot, make(map[int]UserInfo), nil, nil}, nil
}

func NewBotWithCamera(tok string, cam camera.Camera) (*Bot, error) {
	b, err := NewBot(tok)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create bot with camera")
	}
	b.addCamera(cam)
	return b, nil
}

func (b *Bot) addCamera(cam camera.Camera) {
	b.camera = cam
}

func (b *Bot) Setup(bc bot.BotConfig) error {
	for _, usrid := range bc.AllowedUsers {
		uids := strings.Split(usrid, ",")
		for _, uid := range uids {
			cid, err := strconv.Atoi(uid)
			if err != nil {
				return errors.Wrapf(err, "Could not convert user's chat id %s", uid)
			}
			b.allowedUsers[cid] = UserInfo{chatID: int64(cid)}
		}

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
	go b.updateHandler()
	return nil
}

func (b *Bot) Send(m bot.Message) error {
	for _, usrinfo := range b.allowedUsers {
		err := b.handleMessage(m, usrinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) SendTo(usr int, m bot.Message) error {
	usrinfo := b.allowedUsers[usr]
	return b.handleMessage(m, usrinfo)
}

func (b *Bot) handleMessage(m bot.Message, usrinfo UserInfo) error {
	switch m.Format {
	case bot.TEXT:
		msg := tgbotapi.NewMessage(usrinfo.chatID, string(m.Data))
		_, err := b.bot.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "Could not send text message to %d", usrinfo.chatID)
		}
	case bot.IMAGE:
		msg := tgbotapi.NewPhotoUpload(usrinfo.chatID, string(m.Data))
		msg.Caption = string(m.Data)
		_, err := b.bot.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "Could not send photo message to %d", usrinfo.chatID)
		}
	case bot.VIDEO:
		msg := tgbotapi.NewVideoNoteUpload(usrinfo.chatID, 480, string(m.Data))
		msg.Duration = 10
		_, err := b.bot.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "Could not send video note message to %d", usrinfo.chatID)
		}
	}
	return nil
}

func (b *Bot) updateHandler() {
	for u := range b.updates {
		if u.Message == nil { // ignore any non-Message Updates
			continue
		}
		switch u.Message.Text {
		case "/subscribe":
			if _, ok := b.allowedUsers[u.Message.From.ID]; ok {
				b.allowedUsers[u.Message.From.ID] = UserInfo{
					chatID:    u.Message.Chat.ID,
					username:  u.Message.From.UserName,
					firstName: u.Message.From.FirstName,
					lastName:  u.Message.From.LastName,
				}
				greetings := ""
				if u.Message.From.UserName != "" {
					greetings = fmt.Sprintf(helloMessageTemplate, u.Message.From.UserName)
				} else {
					greetings = fmt.Sprintf(helloMessageTemplate, u.Message.From.FirstName+" "+u.Message.From.LastName)
				}
				b.Send(bot.Message{
					Format: bot.TEXT,
					Data:   []byte(greetings),
				})
			}
		case "/picture":
			if _, ok := b.allowedUsers[u.Message.From.ID]; ok {
				filename, err := b.camera.Picture()
				if err != nil {
					b.Send(bot.Message{
						Format: bot.TEXT,
						Data:   []byte(err.Error()),
					})
					return
				}
				b.SendTo(u.Message.From.ID, bot.Message{
					Format: bot.IMAGE,
					Data:   []byte(filename),
				})
			}
		case "/video":
			if _, ok := b.allowedUsers[u.Message.From.ID]; ok {
				filename, err := b.camera.Video()
				if err != nil {
					b.SendTo(u.Message.From.ID, bot.Message{
						Format: bot.TEXT,
						Data:   []byte(err.Error()),
					})
					return
				}
				b.SendTo(u.Message.From.ID, bot.Message{
					Format: bot.VIDEO,
					Data:   []byte(filename),
				})
			}
		}
		log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)
	}
}
