package main

import (
	"os"
	"strings"

	"github.com/indiependente/motionbot/sensor"

	"github.com/indiependente/motionbot/bot"
	"github.com/indiependente/motionbot/bot/telegram"
	"github.com/indiependente/motionbot/sensor/camera"
	"github.com/indiependente/motionbot/sensor/pir"
	log "github.com/sirupsen/logrus"
)

func main() {
	// sensor setup
	pin := os.Getenv("GPIOPIN")
	cfg := pir.SensorConfig{PinOut: pin}
	s := pir.NewSensor(cfg)
	err := s.Setup()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("sensor setup")
	}
	outCh := s.Read()

	// camera setup
	cam := &camera.NoIRCamera{}

	// bot setup
	token := os.Getenv("TOKEN")
	allowedUsers := strings.Split(os.Getenv("ALLOWED_USERS"), ",")
	botcfg := bot.BotConfig{AllowedUsers: allowedUsers}
	tbot, err := telegram.NewBot(token)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("new bot")
	}

	err = tbot.Setup(botcfg)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("bot setup")
	}
	err = tbot.Start()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("bot start")
	}

	for m := range outCh {
		tbot.Send(bot.Message{Format: bot.MessageFormat(m.Format), Data: m.Data})
		filename, err := cam.Picture()
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("picture")
		}
		tbot.Send(bot.Message{Format: sensor.TYPEIMAGE, Data: []byte(filename)})
	}
}
