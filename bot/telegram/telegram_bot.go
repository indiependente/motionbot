package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/indiependente/motionbot/bot"
	"github.com/indiependente/motionbot/sensor/camera"
	"github.com/pkg/errors"
)

const (
	newSubscriberMessageTemplate = "Welcome %s! 😃"
	unsubscribeMessageTemplate   = "Goodbye %s! 👋"
	notAllowedUserTemplate       = "%s You Shall Not Pass! 🧙‍♂️"
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
	activeUsers  map[int]UserInfo
	updates      <-chan tgbotapi.Update
	camera       camera.Camera
}

func NewBot(tok string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(tok)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create bot")
	}
	return &Bot{bot, make(map[int]UserInfo), make(map[int]UserInfo), nil, nil}, nil
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
	for _, usrinfo := range b.activeUsers {
		err := b.handleMessage(m, int(usrinfo.chatID))
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) SendTo(chatID int, m bot.Message) error {
	return b.handleMessage(m, chatID)
}

func (b *Bot) handleMessage(m bot.Message, chatID int) error {
	chatID64 := int64(chatID)
	switch m.Format {
	case bot.TEXT:
		msg := tgbotapi.NewMessage(chatID64, string(m.Data))
		_, err := b.bot.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "Could not send text message to %d", chatID64)
		}
	case bot.IMAGE:
		msg := tgbotapi.NewPhotoUpload(chatID64, string(m.Data))
		msg.Caption = string(m.Data)
		_, err := b.bot.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "Could not send photo message to %d", chatID64)
		}
	case bot.VIDEO:
		msg := tgbotapi.NewVideoNoteUpload(chatID64, 480, string(m.Data))
		msg.Duration = 10
		_, err := b.bot.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "Could not send video note message to %d", chatID64)
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
				b.activeUsers[u.Message.From.ID] = UserInfo{
					chatID:    u.Message.Chat.ID,
					username:  u.Message.From.UserName,
					firstName: u.Message.From.FirstName,
					lastName:  u.Message.From.LastName,
				}
				b.SendTo(
					u.Message.From.ID,
					bot.Message{
						Format: bot.TEXT,
						Data:   []byte(getUsrName(u.Message, newSubscriberMessageTemplate)),
					})
			} else {
				b.SendTo(
					u.Message.From.ID,
					bot.Message{
						Format: bot.TEXT,
						Data:   []byte(getUsrName(u.Message, notAllowedUserTemplate)),
					})
			}
		case "/picture":
			if _, ok := b.allowedUsers[u.Message.From.ID]; ok {
				filename, err := b.camera.Picture()
				if err != nil {
					b.SendTo(
						u.Message.From.ID,
						bot.Message{
							Format: bot.TEXT,
							Data:   []byte(err.Error()),
						})
					return
				}
				b.SendTo(
					u.Message.From.ID,
					bot.Message{
						Format: bot.IMAGE,
						Data:   []byte(filename),
					})
			}
		case "/video":
			if _, ok := b.allowedUsers[u.Message.From.ID]; ok {
				filename, err := b.camera.Video()
				if err != nil {
					b.SendTo(
						u.Message.From.ID,
						bot.Message{
							Format: bot.TEXT,
							Data:   []byte(err.Error()),
						})
					return
				}
				b.SendTo(
					u.Message.From.ID,
					bot.Message{
						Format: bot.VIDEO,
						Data:   []byte(filename),
					})
			}
		case "/unsubscribe":
			if _, ok := b.allowedUsers[u.Message.From.ID]; ok {
				b.SendTo(
					u.Message.From.ID,
					bot.Message{
						Format: bot.TEXT,
						Data:   []byte(getUsrName(u.Message, unsubscribeMessageTemplate)),
					})
				delete(b.activeUsers, u.Message.From.ID)
			}
		}
		log.Printf("[%s] %s", u.Message.From.UserName, u.Message.Text)
	}
}

func getUsrName(m *tgbotapi.Message, template string) string {
	greetings := ""
	if m.From.UserName != "" {
		greetings = fmt.Sprintf(template, m.From.UserName)
	} else {
		greetings = fmt.Sprintf(template, m.From.FirstName+" "+m.From.LastName)
	}
	return greetings
}
