package config

type Config struct {
	App      App
	DB       DB
	Adapters Adapters
}

type App struct {
	ListenAddr string `env:"APP_LISTEN_ADDR" env-default:":8080"`
	Url        string `env:"APP_URL" env-default:"localhost:8080"`
}

type DB struct {
	Dsn string `env:"DB_DSN" env-required:"true"`
}

type Adapters struct {
	Url string `env:"ADAPTERS_URL"`
}
