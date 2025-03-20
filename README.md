# Инструкция

- В `.env` необходимо указать/поменять:
    +POSTGRES_DB=
    +POSTGRES_USER=
    +POSTGRES_PASSWORD=

- В docker-compose также указать
```
      POSTGRES_USER: 
      POSTGRES_PASSWORD: 
      POSTGRES_DATABASE: 
```

в .env-docker:
+ **STORAGE_TYPE** — тип хранилища: `memory` или `postgres`

## Для запуска
`docker compose --profile app up --build -d`

SWAGGER будет находится по пути:

EXAMPLE: `localhost:8080/docs`