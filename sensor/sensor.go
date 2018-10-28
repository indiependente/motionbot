//go:generate mockgen -source=sensor.go -destination sensor_mock.go -package=sensor

package sensor

const (
	// TYPETEXT represents a text Output message format.
	TYPETEXT = "TEXT"
)

// Output is the output sent by the sensor.
type Output struct {
	Format string
	Data   []byte
}

// NewOutput creates an Output containing the format and the data passed as input.
func NewOutput(f string, d []byte) Output {
	return Output{f, d}
}

// Sensor abstracts a sensor.
type Sensor interface {
	Setup() error
	Read() chan Output
}
