package main

import (
	"os"

	"github.com/indiependente/motionbot/sensor"
	"github.com/indiependente/motionbot/sensor/pir"
	log "github.com/sirupsen/logrus"
)

func main() {
	pin := os.Getenv("GPIOPIN")
	cfg := pir.SensorConfig{PinOut: pin}
	s := pir.NewSensor(cfg)
	err := s.Setup()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("setup")
	}
	outCh := s.Read()

	for m := range outCh {
		if m.Format == sensor.TYPETEXT {
			log.WithFields(log.Fields{"type": m.Format, "data": string(m.Data)}).Info("message")
		}
	}
}
