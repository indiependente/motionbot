package pir

import (
	"github.com/indiependente/motionbot/sensor"
	"github.com/pkg/errors"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

// Sensor represents a PIR sensor.
type Sensor struct {
	pin   string
	outCh chan sensor.Output
	p     gpio.PinIO
}

// SensorConfig holds the configuration parameters of a PIR Sensor.
type SensorConfig struct {
	PinIn  string
	PinOut string
}

// NewSensor creates a new sensor that
func NewSensor(conf SensorConfig) Sensor {
	return Sensor{
		pin:   conf.PinOut,
		outCh: make(chan sensor.Output),
	}
}

// Setup initializes the PIR sensor.
func (s *Sensor) Setup() error {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		return errors.Wrap(err, "Could not init host")
	}

	// Lookup a pin by its number:
	s.p = gpioreg.ByName(s.pin)

	// Set it as input.
	if err := s.p.In(gpio.PullNoChange, gpio.RisingEdge); err != nil {
		return errors.Wrap(err, "Could not set pin as input")
	}
	return nil
}

// Read starts a separate goroutine that will write Output on the returned channel.
func (s *Sensor) Read() chan sensor.Output {
	go func() {
		// Wait for edges as detected by the hardware.
		for {
			s.p.WaitForEdge(-1)
			if s.p.Read() == gpio.High {
				s.outCh <- sensor.NewOutput(sensor.TYPETEXT, []byte("Movement detected ⚠️"))
			}
		}
	}()
	return s.outCh
}
