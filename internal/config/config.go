package config

type Config struct {
	App      App
	DB       DB
	Adapters Adapters
}

type App struct {
	Url string `env:"APP_URL" env-default:"localhost:8080"`
}

type DB struct {
	Dsn string `env:"DB_DSN" env-required:"true"`
}

type Adapters struct {
	Url string `env:"ADAPTERS_URL"`
}
