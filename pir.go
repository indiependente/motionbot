package pir

import (
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

const (
	typetext = "TEXT"
)

type Output struct {
	format string
	data   []byte
}

func NewOutput(f string, d []byte) Output {
	return Output{f, d}
}

type Sensor struct {
	pin   string
	outCh chan Output
	p     gpio.PinIO
}

func NewSensor(pin string) Sensor {
	return Sensor{
		pin:   pin,
		outCh: make(chan Output),
	}
}

func (s *Sensor) Setup() error {
	// Load all the drivers:
	if _, err := host.Init(); err != nil {
		return err
	}

	// Lookup a pin by its number:
	p := gpioreg.ByName(s.pin)

	// Set it as input.
	if err := p.In(gpio.PullNoChange, gpio.RisingEdge); err != nil {
		return err
	}
	return nil
}
func (s *Sensor) Read() chan Output {
	go func() {
		// Wait for edges as detected by the hardware.
		for {
			s.p.WaitForEdge(-1)
			if s.p.Read() == gpio.High {
				s.outCh <- Output{format: typetext, data: []byte("Movement detected ⚠️")}

			}
		}

	}()
	return s.outCh
}
