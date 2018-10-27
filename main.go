package main

import (
	"github.com/indiependente/motionbot/sensor"
	"github.com/indiependente/motionbot/sensor/pir"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := pir.SensorConfig{PinOut: "GPIO23"}
	s := pir.NewSensor(cfg)
	err := s.Setup()
	if err != nil {
		log.WithFields(log.Fields{"setup": err}).Fatal("ERROR")
	}
	outCh := s.Read()

	for m := range outCh {
		if m.Format == sensor.TYPETEXT {
			log.WithFields(log.Fields{"type": m.Format, "data": string(m.Data)}).Info("message")
		}
	}
}
