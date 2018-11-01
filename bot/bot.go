package bot

type BotConfig struct {
	AllowedUsers []string
}

type MessageFormat string

const (
	TEXT  MessageFormat = "TEXT"
	IMAGE MessageFormat = "IMAGE"
	VIDEO MessageFormat = "VIDEO"
)

type Message struct {
	Format MessageFormat
	Data   []byte
}

type Bot interface {
	Setup(BotConfig) error
	Start() error
	Send(Message) error
}
