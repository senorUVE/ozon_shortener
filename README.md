# Инструкция

* В `.env` необходимо указать/поменять:
    + POSTGRES_DB=
    + POSTGRES_USER=
    + POSTGRES_PASSWORD=

В **docker-compose** также указать
```
      POSTGRES_USER: 
      POSTGRES_PASSWORD: 
      POSTGRES_DATABASE: 
```

в `.env-docker`:
+ **STORAGE_TYPE** — тип хранилища: `memory` или `postgres`
+ **APP_URL** - если нужно другой домен/порт у коротких ссылок.
+ **DB_DSN** - если меняется хост/порт базы или пароли.

## Для запуска

Требуются env-авры, примеры в файлах


`docker compose --profile app up --build -d`

При запуске создаст все необходимые таблицы в postgre

SWAGGER будет находится по пути как в примере:`localhost:8080/docs` 

Генерация html с покрытием тестов

```bash
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```