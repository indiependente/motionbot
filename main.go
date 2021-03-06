package main

import (
	"os"
	"strings"

	"github.com/indiependente/motionbot/bot"
	"github.com/indiependente/motionbot/bot/telegram"
	"github.com/indiependente/motionbot/sensor/camera"
	"github.com/indiependente/motionbot/sensor/pir"
	"github.com/indiependente/motionbot/video/muxer/ffmpeg"
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
	cam := &camera.NoIRCamera{
		Muxer: ffmpeg.Muxer{Rate: 30},
	}

	// bot setup
	token := os.Getenv("TOKEN")
	allowedUsers := strings.Split(os.Getenv("ALLOWED_USERS"), ",")

	tbot, err := telegram.NewBotWithCamera(token, cam)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("new bot")
	}

	botcfg := bot.Config{AllowedUsers: allowedUsers}
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
	}
}
