package apiserver

type Config struct {
	BindAddr    string
	DatabaseURL string
	SessionKey  string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:   ":5000",
		SessionKey: "jdfhdfdj",
		DatabaseURL: "host=localhost dbname=db-forum_dev sslmode=disable port=5432 password=forum user=ubuntu",
	}
}