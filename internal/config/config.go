package config

type Config struct {
	App App
	DB  DB
}

type App struct {
	ListenAddr  string `env:"APP_LISTEN_ADDR" env-default:":8080"`
	Url         string `env:"APP_URL" env-default:"localhost:8080"`
	StorageType string `env:"STORAGE_TYPE" env-default:"postgres"`
}

type DB struct {
	Dsn string `env:"DB_DSN" env-required:"true"`
}
