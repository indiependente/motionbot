package bot

type BotConfig struct {
	AllowedUsers []string
}

type Message struct {
	Format string
	Data   []byte
}

type Bot interface {
	Setup(BotConfig) error
	Start() error
	Send(Message) error
}
