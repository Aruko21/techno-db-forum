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
		DatabaseURL: "host=localhost dbname=db-forum sslmode=disable port=5432 user=technopark password=park",
	}
}