package bot

// Config represents a bot configuration object.
type Config struct {
	AllowedUsers []string
}

// MessageFormat represents the type of messages supported.
type MessageFormat string

const (
	// TEXT is a text message format.
	TEXT MessageFormat = "TEXT"
	// IMAGE is an image message format.
	IMAGE MessageFormat = "IMAGE"
	// VIDEO is a video message format.
	VIDEO MessageFormat = "VIDEO"
)

// Message represents a message to be sent by the bot.
type Message struct {
	Format MessageFormat
	Data   []byte
}

// Bot abstracts the behavior of a bot.
type Bot interface {
	Setup(Config) error
	Start() error
	Send(Message) error
}
